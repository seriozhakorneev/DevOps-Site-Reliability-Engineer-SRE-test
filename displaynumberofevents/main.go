package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	filePath     = "displaynumberofevents/events.log"
	eventSuffix  = "NOK"
	parseLayout  = "[2006-01-02 15:04:05]"
	outputLayout = "[2006-01-02 15:04]"
)

type eventCount struct {
	minuteTime time.Time
	count      int
}

func main() {
	logFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file path(%s): %s", filePath, err)
	}
	defer logFile.Close()

	result, err := getEventsInMinute(
		bufio.NewScanner(logFile),
		parseLayout,
		eventSuffix,
	)
	if err != nil {
		log.Fatalf("Get events in minute failed: %s", err)
	}

	for _, v := range result[1:] {
		fmt.Println(v.minuteTime.Format(outputLayout), v.count)
	}
}

// getEventsInMinute
// receives *bufio.Scanner, made it just for speed, don't want to make +1 loop storing lines to slice.
func getEventsInMinute(scanner *bufio.Scanner, layout, suffix string) ([]eventCount, error) {
	events := make([]eventCount, 1)

	for scanner.Scan() {
		// skip empty line
		if len(scanner.Text()) < len(layout) {
			continue
		}

		//	parse string date to golang time, truncate to minute
		//  also stop runtime and print error if line did not validate format
		t, err := parseTimeTo(scanner.Text()[:len(layout)], layout, time.Minute)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse time string(%s), layout(%s), error: %w",
				scanner.Text()[:len(layout)],
				layout,
				err,
			)
		}

		if events[len(events)-1].minuteTime != t {
			events = append(
				events,
				eventCount{
					minuteTime: t,
					count:      0,
				},
			)
		}

		if strings.HasSuffix(scanner.Text(), suffix) {
			events[len(events)-1].count++
		}
	}

	return events, nil
}

func parseTimeTo(s, l string, d time.Duration) (time.Time, error) {
	t, err := time.Parse(l, s)
	if err != nil {
		return time.Time{}, err
	}

	return t.Truncate(d), nil
}
