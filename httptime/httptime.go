package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"sync"
	"sync/atomic"
	"time"
)

var (
	start, connect, dns, tlsHandshake time.Time
)

var (
	timeStart = time.Now()
	mutex     = &sync.Mutex{}

	requests int64
	avg      int64
)

var (
	numReq = flag.Int("n", 4, "Number of requets")
	url    = flag.String("url", "https://www.website.tld/", "Target URL for Testing")
)

func getHTTPTime(url string, wg *sync.WaitGroup) int64 {
	defer wg.Done()

	req, _ := http.NewRequest("GET", url, nil)

	trace := &httptrace.ClientTrace{
		DNSStart:             func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone:              func(ddi httptrace.DNSDoneInfo) {},
		TLSHandshakeStart:    func() { tlsHandshake = time.Now() },
		TLSHandshakeDone:     func(cs tls.ConnectionState, err error) {},
		ConnectStart:         func(network, addr string) { connect = time.Now() },
		ConnectDone:          func(network, addr string, err error) {},
		GotFirstResponseByte: func() {},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatal(err)
	}

	atomic.AddInt64(&requests, 1)
	return (time.Since(timeStart).Milliseconds())
}

func main() {
	var wg sync.WaitGroup

	flag.Parse()

	for i := 0; i < *numReq; i++ {
		wg.Add(1)
		go func() {
			var result = getHTTPTime(*url, &wg)
			mutex.Lock()
			avg = ((avg + result) / requests)
			mutex.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println(avg)
}
