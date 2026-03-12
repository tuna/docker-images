package main

import (
	"flag"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type config struct {
	portNumber   uint16
	queueSize    int
	updatePeriod time.Duration
}

func parseCLIParams() config {
	var cfg config
	var portNumber int

	flag.IntVar(&portNumber, "port-number", 8888, "Port number to listen on")
	flag.IntVar(&cfg.queueSize, "queue-size", 100, "Queue capacity")
	flag.DurationVar(&cfg.updatePeriod, "update-period", 1*time.Second, "Queue update period")
	flag.Parse()

	if portNumber < 1 || portNumber > 65535 {
		panic("port-number must be between 1 and 65535")
	}
	cfg.portNumber = uint16(portNumber)

	if cfg.queueSize <= 0 {
		panic("queue-size must be greater than 0")
	}
	if cfg.updatePeriod <= 0 {
		panic("update-period must be greater than 0")
	}

	return cfg
}

var idCounter uint64 = 0

func generateID() uint64 {
	if idCounter == ^uint64(0) {
		panic("ID counter overflow")
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

	clientQueue := make(chan chan struct{})

	for i := 0; i < cfg.queueSize; i++ {
		go func() {
			for {
				done := <-clientQueue
				<-done
			}
		}()
	}

	connections := make(chan net.Conn)

	listener := func(addr string) {
		connListener, err := net.Listen("tcp", addr)
		if err != nil {
			panic(err)
		}
		defer connListener.Close()
		for {
			conn, err := connListener.Accept()
			if err != nil {
				fmt.Printf("Failed to accept connection: %v\n", err)
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
