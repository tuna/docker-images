package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type config struct {
	listeningAddress     string
	upstreamAddress      string
	maxResponseDeadline  time.Duration
	maxBuffers           uint64
	bufferUnitSize       uint64
	onDiskBufferPath     string
	onDiskBufferUnitSize uint64
}

type UnitValue uint64

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

func parseCLIParams() config {

	// Define command-line flags
	listeningAddress := flag.String("listen", ":8080", "Address to listen on (e.g., :8080)")
	upstreamAddress := flag.String("upstream", "localhost:80", "Upstream server address (e.g., localhost:80)")
	maxResponseDeadline := flag.Duration("max-deadline", 1*time.Second, "Maximum response deadline (e.g., 1s)")
	var maxBufferSize UnitValue = 1 * 1024 * 1024 * 1024 // Default to 1GB
	flag.Var(&maxBufferSize, "max-buffer-size", "Maximum buffer size in memory")
	var bufferUnitSize UnitValue = UnitValue(os.Getpagesize()) // Default to system page size
	flag.Var(&bufferUnitSize, "buffer-unit-size", "Buffer unit size in memory")
	var onDiskBufferUnitSize UnitValue = 4 * 1024 * 1024 // Default to 4MB
	flag.Var(&onDiskBufferUnitSize, "on-disk-buffer-unit-size", "Buffer unit size on disk")
	onDiskBufferPath := flag.String("on-disk-buffer-path", "./buffer", "Path for on-disk buffer")

	// Parse command-line flags
	flag.Parse()
	// Validate CLI parameters that must be positive to avoid panics or busy loops.
	if bufferUnitSize == 0 {
		fmt.Fprintln(os.Stderr, "invalid value for --buffer-unit-size: must be greater than zero")
		os.Exit(2)
	}
	if onDiskBufferUnitSize == 0 {
		fmt.Fprintln(os.Stderr, "invalid value for --on-disk-buffer-unit-size: must be greater than zero")
		os.Exit(2)
	}
	if *maxResponseDeadline <= 0 {
		fmt.Fprintln(os.Stderr, "invalid value for --max-deadline: must be greater than zero")
		os.Exit(2)
	}
	maxBuffers := maxBufferSize / bufferUnitSize
	if maxBuffers == 0 {
		panic("max-buffer-size must be at least as large as buffer-unit-size")
	}

	return config{
		listeningAddress:     *listeningAddress,
		upstreamAddress:      *upstreamAddress,
		maxResponseDeadline:  *maxResponseDeadline,
		maxBuffers:           uint64(maxBuffers),
		bufferUnitSize:       uint64(bufferUnitSize),
		onDiskBufferPath:     *onDiskBufferPath,
		onDiskBufferUnitSize: uint64(onDiskBufferUnitSize),
	}
}

type Buffer struct {
	memBuf  *[]byte
	diskBuf *os.File
}

var NrInMemBuffers atomic.Uint64

type StreamConn interface {
	net.Conn
	CloseRead() error
	CloseWrite() error
}

func handleConnection(conn *net.TCPConn, cfg config) {
	var _uc net.Conn
	var err error
	if strings.HasPrefix(cfg.upstreamAddress, "unix:") {
		socketPath := strings.TrimPrefix(cfg.upstreamAddress, "unix:")
		_uc, err = net.Dial("unix", socketPath)
	} else {
		_uc, err = net.Dial("tcp", cfg.upstreamAddress)
	}
	if err != nil {
		fmt.Printf("Failed to connect to upstream server: %v\n", err)
		conn.Close()
		return
	}
	upstreamConn, ok := _uc.(StreamConn)
	if !ok {
		fmt.Printf("Failed to cast upstream connection to TCPConn\n")
		conn.Close()
		_uc.Close()
		return
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() { //Uploading data from client to upstream, no buffering
		defer wg.Done()
		_, err := io.Copy(upstreamConn, conn)
		if err != nil {
			fmt.Printf("Error copying from client to upstream: %v\n", err)
		}
		upstreamConn.CloseWrite()
		conn.CloseRead()
	}()

	go func() { //Downloading data from upstream to client, with buffering
		channel := make(chan Buffer, int(cfg.maxBuffers))
		go func() { // Sender
			defer wg.Done()
			closed := false
			handleSendError := func() {
				upstreamConn.CloseRead()
			}
			for buf := range channel {
				if buf.memBuf != nil {
					if !closed {
						_, err := io.Copy(conn, bytes.NewReader(*buf.memBuf))
						if err != nil {
							fmt.Printf("Error writing to client: %v\n", err)
							closed = true
							conn.CloseWrite()
							handleSendError()
						}
					}
					NrInMemBuffers.Add(^uint64(0))
				} else {
					if !closed {
						_, err := io.Copy(conn, buf.diskBuf)
						if err != nil {
							fmt.Printf("Error writing to client from disk buffer: %v\n", err)
							closed = true
							conn.CloseWrite()
							handleSendError()
						}
					}
					buf.diskBuf.Close()
				}
			}
			if !closed {
				conn.CloseWrite()
			}
		}()
		go func() { // Receiver
			defer wg.Done()
			for {
				upstreamConn.SetReadDeadline(time.Now().Add(cfg.maxResponseDeadline))
				var inMem bool
				for {
					curNrBuffers := NrInMemBuffers.Load()
					if curNrBuffers < cfg.maxBuffers {
						ok := NrInMemBuffers.CompareAndSwap(curNrBuffers, curNrBuffers+1)
						if ok {
							inMem = true
							break
						}
					} else {
						inMem = false
						break
					}
				}
				var readErr error
				if inMem {
					buf := make([]byte, int(cfg.bufferUnitSize))
					var n int
					n, readErr = io.ReadFull(upstreamConn, buf)
					if n > 0 {
						buf = buf[:n]
						channel <- Buffer{memBuf: &buf}
					} else {
						NrInMemBuffers.Add(^uint64(0))
					}
				} else {
					tmpFile, err := os.CreateTemp(cfg.onDiskBufferPath, "buffer-*")
					if err != nil {
						fmt.Printf("Error creating temp file for disk buffer: %v\n", err)
						close(channel)
						upstreamConn.CloseRead()
						break
					}
					err = os.Remove(tmpFile.Name())
					if err != nil {
						fmt.Printf("Error removing temp file after creation: %v\n", err)
						tmpFile.Close()
						close(channel)
						upstreamConn.CloseRead()
						break
					}
					var n int64
					n, readErr = io.CopyN(tmpFile, upstreamConn, int64(cfg.onDiskBufferUnitSize))
					_, err = tmpFile.Seek(0, io.SeekStart)
					if err != nil {
						fmt.Printf("Error seeking to start of temp file: %v\n", err)
						tmpFile.Close()
						close(channel)
						upstreamConn.CloseRead()
						break
					}
					if n > 0 {
						channel <- Buffer{diskBuf: tmpFile}
					} else {
						tmpFile.Close()
					}
				}
				if readErr != nil {
					if ne, ok := readErr.(net.Error); ok && ne.Timeout() {
						continue
					}
					if readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
						fmt.Printf("Error reading from upstream: %v\n", readErr)
					}
					close(channel)
					upstreamConn.CloseRead()
					break
				}
			}
		}()
	}()
	go func() {
		wg.Wait()
		upstreamConn.Close()
		conn.Close()
	}()
}

func main() {
	cfg := parseCLIParams()

	if err := os.MkdirAll(cfg.onDiskBufferPath, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating on-disk buffer directory %q: %v\n", cfg.onDiskBufferPath, err)
		os.Exit(1)
	}

	if cfg.maxBuffers > uint64(int(^uint(0)>>1)) {
		fmt.Fprintf(os.Stderr, "maxBuffers value %d exceeds maximum allowed\n", cfg.maxBuffers)
		os.Exit(1)
	}

	if cfg.bufferUnitSize > uint64(int(^uint(0)>>1)) {
		fmt.Fprintf(os.Stderr, "bufferUnitSize value %d exceeds maximum allowed\n", cfg.bufferUnitSize)
		os.Exit(1)
	}

	if cfg.onDiskBufferUnitSize > ((^uint64(0)) >> 1) {
		fmt.Fprintf(os.Stderr, "onDiskBufferUnitSize value %d exceeds maximum allowed\n", cfg.onDiskBufferUnitSize)
		os.Exit(1)
	}

	connListener, err := net.Listen("tcp", cfg.listeningAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting listener: %v\n", err)
		os.Exit(1)
	}
	defer connListener.Close()
	fmt.Printf("Listening on %s and forwarding to %s\n", cfg.listeningAddress, cfg.upstreamAddress)

	for {
		conn, err := connListener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			go handleConnection(tcpConn, cfg)
		} else {
			fmt.Printf("Failed to cast connection to TCPConn\n")
			conn.Close()
		}
	}
}
