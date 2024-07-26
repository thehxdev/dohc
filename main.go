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
)

func main() {
	err := os.Setenv("GOGC", "50")
	if err != nil {
		log.Fatal(err)
	}

	configureCmdFlags()

	fd, err := os.Open(serversFile)
	if err != nil {
		log.Fatal(err)
	}

	limitCh := make(chan bool, limit)
	scanner := bufio.NewScanner(fd)
	c := createHttpClient()
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	workingDomains := make([]string, 0)

	for scanner.Scan() {
		wg.Add(1)
		line := strings.TrimSpace(scanner.Text())

		go func(domain string) {
			defer wg.Done()

			limitCh <- false
			defer deque(limitCh)

			req, err := http.NewRequestWithContext(context.Background(), "POST", line, createDNSPacket())
			if err != nil {
				log.Println(err)
				return
			}
			req.Header.Set("Content-Type", "application/dns-message")
			req.Header.Set("Accept", "application/dns-message")

			resp, err := c.Do(req)
			if err != nil {
				log.Printf("[FAIL] %s\n", domain)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				log.Printf("[OK] %s\n", domain)
				mu.Lock()
				workingDomains = append(workingDomains, domain)
				mu.Unlock()
			} else {
				log.Printf("[FAIL] %s\n", domain)
			}
		}(line)
	}
	wg.Wait()

	fmt.Printf("%+v\n\n", workingDomains)

	outFileName := "results_" + strings.ReplaceAll(string(time.Now().Format(time.DateTime)), " ", "_") + "_.txt"
	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range workingDomains {
		outFile.Write([]byte(server + "\n"))
	}

	log.Println("wrote results to", outFileName)
}

func configureCmdFlags() {
	flag.StringVar(&serversFile, "f", "doh_servers.txt", "path to servers list file")
	flag.StringVar(&localResolver, "r", "9.9.9.9", "local DNS resolver (to resolve DoH addresses)")
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
		Timeout: time.Second * 10,
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
		ANCount: 1,
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
	inBytes := serializeBuff.Bytes()

	buff := &bytes.Buffer{}
	buff.Write(inBytes)
	return buff
}

func deque(ch <-chan bool) {
	<-ch
}
