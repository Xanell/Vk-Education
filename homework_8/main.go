package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	queryCount := flag.Int("k", 5, "кол-во запросов")
	searchValue := flag.String("q", "Go", "искомая строка")
	flag.Parse()

	var total atomic.Int64
	wg := sync.WaitGroup{}
	ch := make(chan string)
	pool := make(chan struct{}, *queryCount)

	wg.Go(func() {
		reader(ch)
	})

	for url := range ch {
		pool <- struct{}{}
		wg.Go(func() {
			defer func() { <-pool }()
			worker(url, *searchValue, &total)
		})
	}

	wg.Wait()
	fmt.Printf("Total %d", total.Load())
}

// писатель
func reader(ch chan<- string) {
	defer close(ch)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, `\n`)
		for _, rawURL := range parts {
			rawURL = strings.TrimSpace(rawURL)
			if rawURL == "" {
				continue
			}
			ch <- rawURL
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при чтении stdin:", err)
		os.Exit(1)
	}
}

// читатель
func worker(url string, searchValue string, total *atomic.Int64) {
	count, err := processURL(url, searchValue)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("Count for %s: %d\n", url, count)
	total.Add(int64(count))
}

func processURL(url, searchValue string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("сайт %s не доступен: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения ответа %s: %w", url, err)
	}

	count := strings.Count(string(body), searchValue)
	return count, nil
}
