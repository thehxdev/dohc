package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	limit         int
	serversFile   string
	localResolver string
	outFileName   string
)

func init() {
	err := os.Setenv("GOGC", "50")
	if err != nil {
		log.Fatal(err)
	}
}

type Result struct {
	time.Duration
	addr string
}

func main() {
	configureCmdFlags()

	fd, err := os.Open(serversFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	limitCh := make(chan struct{}, limit)
	scanner := bufio.NewScanner(fd)
	c := createHttpClient()
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	workingDomains := make([]Result, 0)

	for scanner.Scan() {
		wg.Add(1)
		line := strings.TrimSpace(scanner.Text())

		go func(dohAddr string) {
			defer wg.Done()

			limitCh <- struct{}{}
			defer func(){
				<-limitCh
			}()

			start := time.Now()

			req, err := http.NewRequestWithContext(context.Background(), "POST", line, createDNSPacket())
			if err != nil {
				log.Println(err)
				return
			}
			req.Header.Set("Content-Type", "application/dns-message")
			req.Header.Set("Accept", "application/dns-message")

			resp, err := c.Do(req)
			if err != nil {
				return
			}
			deltaTime := time.Now().Sub(start)
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				log.Printf("[OK] %4d - %s\n", deltaTime.Milliseconds(), dohAddr)
				mu.Lock()
				workingDomains = append(workingDomains, Result{Duration: deltaTime, addr: dohAddr})
				mu.Unlock()
			}
		}(line)
	}
	wg.Wait()

	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range workingDomains {
		outFile.Write(fmt.Appendf([]byte{}, "%4d %s\n", r.Duration.Milliseconds(), r.addr))
	}

	log.Println("wrote results to", outFileName)
}

func configureCmdFlags() {
	flag.StringVar(&serversFile, "f", "doh_servers.txt", "path to servers list file")
	flag.StringVar(&localResolver, "r", "9.9.9.9", "local DNS resolver (to resolve DoH domain names)")
	flag.StringVar(&outFileName, "o", "results.txt", "path to save scan results")
	flag.IntVar(&limit, "l", max(1, runtime.NumCPU()), "this number of doh servers will be checked concurrently")
	flag.Parse()
}

func createHttpClient() *http.Client {
	dialer := &net.Dialer{
		Timeout: time.Second * 5,
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial(network, net.JoinHostPort(localResolver, "53"))
			},
		},
	}

	return &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
	}
}

func createDNSPacket() *bytes.Buffer {
	dns := layers.DNS{
		ID:      uint16(rand.Intn(0x10000)),
		QR:      true,
		OpCode:  0,
		QDCount: 1,
		ANCount: 0,
		RD:      true,
		RA:      true,
		Questions: []layers.DNSQuestion{
			{
				Name:  []byte("www.google.com"),
				Type:  layers.DNSTypeA,
				Class: layers.DNSClassIN,
			},
		},
	}

	serializeBuff := gopacket.NewSerializeBuffer()
	dns.SerializeTo(serializeBuff, gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: false})

	buff := &bytes.Buffer{}
	buff.Write(serializeBuff.Bytes())
	return buff
}
