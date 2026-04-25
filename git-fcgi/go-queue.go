package main

import (
	goflag "flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pingcap/tidb/v8/pkg/util/cgroup"
	"github.com/prometheus/procfs"

	flag "github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

type UnitValue uint64

func (v *UnitValue) Type() string {
	return "UnitValue"
}

func (v *UnitValue) String() string {
	return strconv.FormatUint(uint64(*v), 10)
}

func (v *UnitValue) Set(s string) error {
	s = strings.TrimSpace(s)
	val := strings.ToUpper(s)

	var multiplier uint64 = 1
	switch {
	case strings.HasSuffix(val, "K"):
		multiplier = 1024
		val = strings.TrimSuffix(val, "K")
	case strings.HasSuffix(val, "M"):
		multiplier = 1024 * 1024
		val = strings.TrimSuffix(val, "M")
	case strings.HasSuffix(val, "G"):
		multiplier = 1024 * 1024 * 1024
		val = strings.TrimSuffix(val, "G")
	case strings.HasSuffix(val, "T"):
		multiplier = 1024 * 1024 * 1024 * 1024
		val = strings.TrimSuffix(val, "T")
	}

	value, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number format: %v", err)
	}
	*v = UnitValue(value * multiplier)
	return nil
}

type config struct {
	portNumber      uint16
	queueSize       int
	updatePeriod    time.Duration
	memoryWatermark uint64
}

func parseCLIParams() config {
	var cfg config
	var portNumber int
	var memoryWatermark UnitValue

	logFlags := goflag.NewFlagSet("logging", goflag.ExitOnError)
	klog.InitFlags(logFlags)
	flag.IntVar(&portNumber, "port-number", 8888, "Port number to listen on")
	flag.IntVar(&cfg.queueSize, "queue-size", 100, "Queue capacity")
	flag.Var(&memoryWatermark, "memory-watermark", "Memory watermark")
	flag.DurationVar(&cfg.updatePeriod, "update-period", 1*time.Second, "Queue update period")
	flag.CommandLine.AddGoFlagSet(logFlags)
	flag.CommandLine.SortFlags = false
	flag.Parse()

	if portNumber < 1 || portNumber > 65535 {
		klog.Exitf("Invalid port number: %d. Must be between 1 and 65535.", portNumber)
	}
	cfg.portNumber = uint16(portNumber)

	if cfg.queueSize <= 0 {
		klog.Exitf("Invalid queue size: %d. Must be greater than 0.", cfg.queueSize)
	}
	if cfg.updatePeriod <= 0 {
		klog.Exitf("Invalid update period: %v. Must be greater than 0.", cfg.updatePeriod)
	}
	cfg.memoryWatermark = uint64(memoryWatermark)
	return cfg
}

var idCounter uint64 = 0

func generateID() uint64 {
	if idCounter == ^uint64(0) {
		klog.Fatal("ID counter overflow")
	}
	idCounter++
	return idCounter
}

var currentID atomic.Uint64

func getCurrentID() uint64 {
	return currentID.Load()
}

func updateCurrentID(newID uint64) {
	for {
		oldID := currentID.Load()
		if oldID >= newID {
			return
		}
		if currentID.CompareAndSwap(oldID, newID) {
			return
		}
	}
}

func checkMemoryLimit(watermark uint64) (bool, uint64, error) {
	pf, err := procfs.NewDefaultFS()
	if err != nil {
		return false, 0, err
	}
	meminfo, err := pf.Meminfo()
	if err != nil {
		return false, 0, err
	}
	if meminfo.MemAvailable == nil {
		return false, 0, fmt.Errorf("MemAvailable is not available in /proc/meminfo")
	}
	systemAvailable := *meminfo.MemAvailable * 1024
	klog.V(2).Infof("/proc/meminfo reports: %d", systemAvailable)
	cgroupLimit, err := cgroup.GetMemoryLimit()
	if err != nil {
		return false, 0, err
	}
	cgroupUsed, err := cgroup.GetMemoryUsage()
	if err != nil {
		return false, 0, err
	}
	cgroupAvailable := cgroupLimit - cgroupUsed
	klog.V(2).Infof("cgroup reports: %d", cgroupAvailable)
	minAvailable := min(cgroupAvailable, systemAvailable)

	return minAvailable >= watermark, minAvailable, nil
}

type memoryChecker struct {
	watermark uint64
	waiter    atomic.Pointer[chan struct{}]
}

func newMemoryChecker(watermark uint64) *memoryChecker {
	return &memoryChecker{
		watermark: watermark,
	}
}

func (mc *memoryChecker) waitForMemoryAvailable() {
	if mc.watermark == 0 {
		return
	}
	var waitCh *chan struct{}
	for {
		if ch := mc.waiter.Load(); ch != nil {
			<-*ch
			return
		}
		newCh := make(chan struct{})
		waitCh = &newCh
		if mc.waiter.CompareAndSwap(nil, waitCh) {
			break
		}
	}
	go func() {
		tiker := time.NewTicker(1 * time.Second)
		defer tiker.Stop()
		firstCheck := true
		for {
			limitOk, available, err := checkMemoryLimit(mc.watermark)
			if err != nil {
				klog.Errorf("Failed to check memory limit: %v\n", err)
				break
			}
			if limitOk {
				if !firstCheck {
					klog.Infof("Memory becomes available %d bytes, above the watermark %d bytes\n", available, mc.watermark)
				}
				break
			}
			if firstCheck {
				klog.Infof("Memory available %d bytes is below the watermark %d bytes\n", available, mc.watermark)
				firstCheck = false
			}
			<-tiker.C
		}
		mc.waiter.Store(nil)
		close(*waitCh)
	}()
	<-*waitCh
}

func servClient(conn net.Conn, waitingQueue chan chan struct{}, reportInterval time.Duration) {

	myId := generateID()

	wakeUp := make(chan struct{})
	closed := make(chan struct{})

	go func() { // Join the queue and wait for my turn
		select {
		case waitingQueue <- closed:
			close(wakeUp)
		case <-closed:
		}
	}()

	go func() { // Watch for client disconnection
		for {
			readBuf := make([]byte, 1)
			_, err := conn.Read(readBuf)
			if err != nil {
				// Client disconnected
				close(closed)
				conn.Close()
				return
			}
		}
	}()

	go func() { // Handle communication with the client
		timer := time.NewTimer(reportInterval)
		defer timer.Stop()
	waitForWakeUp:
		for {
			select {
			case <-wakeUp:
				// I'm being served, update the current ID
				updateCurrentID(myId)
				conn.Write([]byte(fmt.Sprintf("%d\n", 0)))
				break waitForWakeUp
			case <-closed:
				break waitForWakeUp
			case <-timer.C:
				// Report my position in the queue
				current := getCurrentID()
				position := uint64(1)
				if current < myId {
					position = myId - current
				}
				conn.Write([]byte(fmt.Sprintf("%d\n", position)))
				timer.Reset(reportInterval)
			}
		}
	}()
}

func main() {
	cfg := parseCLIParams()

	if cfg.memoryWatermark > 0 {
		_, _, err := checkMemoryLimit(0)
		if err != nil {
			klog.Fatalf("Cannot check memory limit: %v", err)
		}
	}

	clientQueue := make(chan chan struct{})

	memoryChecker := newMemoryChecker(cfg.memoryWatermark)

	for i := 0; i < cfg.queueSize; i++ {
		go func() {
			for {
				memoryChecker.waitForMemoryAvailable()
				done := <-clientQueue
				<-done
			}
		}()
	}

	connections := make(chan net.Conn)

	listener := func(addr string) {
		connListener, err := net.Listen("tcp", addr)
		if err != nil {
			klog.Exitf("Failed to listen on %s: %v", addr, err)
		}
		defer connListener.Close()
		for {
			conn, err := connListener.Accept()
			if err != nil {
				klog.Errorf("Failed to accept connection: %v", err)
				continue
			}
			connections <- conn
		}
	}
	go listener(fmt.Sprintf("127.0.0.1:%d", cfg.portNumber))
	go listener(fmt.Sprintf("[::1]:%d", cfg.portNumber))

	for {
		conn := <-connections
		servClient(conn, clientQueue, cfg.updatePeriod)
	}
}
