package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	url      = flag.String("url", "http://localhost:8080", "target URL")
	duration = flag.Duration("duration", 1*time.Minute, "test duration")
	workers  = flag.Int("workers", 10, "number of workers")
	requests = flag.Int("requests", 1000, "requests per worker")
)

type result struct {
	duration  time.Duration
	status    int
	error     error
	bytesSent int64
}

type statistics struct {
	total     int
	success   int
	failed    int
	durations []time.Duration
	bytes     int64
	mu        sync.Mutex
}

func main() {
	flag.Parse()

	results := make(chan *result, *workers**requests)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(i, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var stats statistics
	for r := range results {
		stats.add(r)
	}

	stats.print(time.Since(start))
}

func (s *statistics) add(r *result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.total++
	if r.error == nil && r.status == http.StatusOK {
		s.success++
	} else {
		s.failed++
	}
	s.durations = append(s.durations, r.duration)
	s.bytes += r.bytesSent
}

func (s *statistics) print(duration time.Duration) {
	fmt.Printf("\nTest Results:\n")
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Total Requests: %d\n", s.total)
	fmt.Printf("Successful Requests: %d\n", s.success)
	fmt.Printf("Failed Requests: %d\n", s.failed)
	fmt.Printf("Requests/sec: %.2f\n", float64(s.total)/duration.Seconds())
	fmt.Printf("Transfer/sec: %.2f MB\n", float64(s.bytes)/(1024*1024)/duration.Seconds())

	if len(s.durations) > 0 {
		var total time.Duration
		for _, d := range s.durations {
			total += d
		}
		avg := total / time.Duration(len(s.durations))
		fmt.Printf("Average Response Time: %v\n", avg)
	}
}

func worker(id int, wg *sync.WaitGroup, results chan<- *result) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for i := 0; i < *requests; i++ {
		start := time.Now()
		resp, err := client.Get(*url)
		r := &result{
			duration: time.Since(start),
		}

		if err != nil {
			r.error = err
		} else {
			r.status = resp.StatusCode
			r.bytesSent = resp.ContentLength
			resp.Body.Close()
		}

		results <- r
	}
}
