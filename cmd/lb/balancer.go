package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"hash/adler32"
	"strconv"
	"github.com/SunRiseGG/ArchitectureLab3/httptools"
	"github.com/SunRiseGG/ArchitectureLab3/signal"
)

var (
	port = flag.Int("port", 8090, "load balancer port")
	timeoutSec = flag.Int("timeout-sec", 3, "request timeout time in seconds")
	https = flag.Bool("https", false, "whether backends support HTTPs")

	traceEnabled = flag.Bool("trace", false, "whether to include tracing information into responses")
)

var (
	timeout = time.Duration(*timeoutSec) * time.Second
	serversPool = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}
	healthyServers []int
)

func scheme() string {
	if *https {
		return "https"
	}
	return "http"
}

func contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func removeByValue(slice []int, key int) []int {
	isExist := contains(slice, key)
  if isExist {
		valueIndex := -1
		for i, n := range slice {
			if key == n {
				valueIndex = i
			}
		}
    return append(slice[:valueIndex], slice[valueIndex+1:]...)
  }
	return slice
}

func createHash(str string) int {
		hashValue := adler32.Checksum([]byte(str))
		result, _ := strconv.Atoi(strconv.FormatUint(uint64(hashValue), 10))
		return result
}

func checkServer(server string, key int) {
	healthy := health(server)
	log.Println(server, healthy)
	if(healthy) {
		if !contains(healthyServers, key) {
			healthyServers = append(healthyServers, key)
		}
	} else {
		healthyServers = removeByValue(healthyServers, key)
	}
}

func chooseServer(addr string) string {
	hash := createHash(addr)
	if contains(healthyServers, hash % len(serversPool)) {
		return serversPool[hash % len(serversPool)]
	} else {
		return serversPool[healthyServers[hash % len(healthyServers)]]
	}
}

func health(dst string) bool {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s://%s/health", scheme(), dst), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func forward(dst string, rw http.ResponseWriter, r *http.Request) error {
	ctx, _ := context.WithTimeout(r.Context(), timeout)
	fwdRequest := r.Clone(ctx)
	fwdRequest.RequestURI = ""
	fwdRequest.URL.Host = dst
	fwdRequest.URL.Scheme = scheme()
	fwdRequest.Host = dst

	resp, err := http.DefaultClient.Do(fwdRequest)
	if err == nil {
		for k, values := range resp.Header {
			for _, value := range values {
				rw.Header().Add(k, value)
			}
		}
		if *traceEnabled {
			rw.Header().Set("lb-from", dst)
		}
		log.Println("fwd", resp.StatusCode, resp.Request.URL)
		rw.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		_, err := io.Copy(rw, resp.Body)
		if err != nil {
			log.Printf("Failed to write response: %s", err)
		}
		return nil
	} else {
		log.Printf("Failed to get response from %s: %s", dst, err)
		rw.WriteHeader(http.StatusServiceUnavailable)
		return err
	}
}

func main() {
	flag.Parse()

	checkServer("server1:8080", 0)
	checkServer("server2:8080", 1)
	checkServer("server3:8080", 2)

		// TODO: Використовуйте дані про стан сервреа, щоб підтримувати список тих серверів, яким можна відправляти ззапит.
	for key, server := range serversPool {
		key := key
		server := server
		go func() {
			for range time.Tick(10 * time.Second) {
				checkServer(server, key)
			}
		}()
	}

	frontend := httptools.CreateServer(*port, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// TODO: Рееалізуйте свій алгоритм балансувальника.
		if (len(healthyServers) != 0) {
		  forward(chooseServer(req.URL.Path), rw, req)
		} else {
		  log.Println("All servers are busy. Wait please.")
	}
	}))

	log.Println("Starting load balancer...")
	log.Printf("Tracing support enabled: %t", *traceEnabled)
	frontend.Start()
	signal.WaitForTerminationSignal()
}
