package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	frequency    = 1 * time.Minute
	httpPrefix   = "http://"
	metricPath   = "/api/count"
	outputLayout = "2006-01-02 15:04:00"
)

var servers = [...]string{"maria.ru", "rose.ru", "sina.ru"}

type output struct {
	Count *int `json:"count"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// starts after minute of waiting
	tick := time.Tick(frequency)
	for {
		select {
		case s := <-interrupt:
			log.Fatal("Received signal:", s.String())
		case x := <-tick:
			for _, server := range servers {
				count, err := getCount(httpPrefix + server + metricPath)
				if err != nil {
					fmt.Println(x.Format(outputLayout), server, err)
				} else {
					fmt.Println(x.Format(outputLayout), server, count)
				}
			}
		}
	}
}

func getCount(path string) (int, error) {
	response, err := http.Get(path)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(
			"response status code is not %d: Status Code: %d",
			http.StatusOK,
			response.StatusCode,
		)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		return 0, fmt.Errorf(
			"content-type header is not application/json: Content-Type: %s",
			response.Header.Get("Content-Type"),
		)
	}

	o := output{}
	err = json.NewDecoder(response.Body).Decode(&o)
	if err != nil {
		return 0, fmt.Errorf("decode json failed: %w", err)
	}

	if o.Count == nil {
		return 0, fmt.Errorf("response data is empty")
	}

	return *o.Count, nil
}
