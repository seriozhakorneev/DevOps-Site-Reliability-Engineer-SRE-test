package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

const (
	tBrokenFilePath = "testfiles/1.log"
	tFilePath       = "testfiles/2.log"
)

func TestParseTimeToError(t *testing.T) {
	tString := "2006-01-02 15:04:05"
	tLayout := "[2006-01-02 15:04:05]"
	tDuration := time.Minute

	expError := fmt.Errorf("parsing time \"2006-01-02 15:04:05\" as" +
		" \"[2006-01-02 15:04:05]\":" +
		" cannot parse \"2006-01-02 15:04:05\" as \"[\"")

	resTime, err := parseTimeTo(tString, tLayout, tDuration)

	if !resTime.IsZero() {
		t.Fatalf("Expected zero time in testing: %s, Got: %s", time.Time{}, resTime)
	}

	if err.Error() != expError.Error() {
		t.Fatalf("Expected error: %s, Got: %s", expError, err)
	}
}

func TestParseTimeToResult(t *testing.T) {
	tString := "[2006-01-02 15:04:05]"
	tInLayout := "[2006-01-02 15:04:05]"
	tOutLayout := "[2006-01-02 15]"
	tDuration := time.Hour

	expResult := "[2006-01-02 15]"
	resTime, err := parseTimeTo(tString, tInLayout, tDuration)
	if err != nil {
		t.Fatalf("Unexpected error in test: %s", err)
	}

	format := resTime.Format(tOutLayout)
	if expResult != format {
		t.Fatalf("Expected result: %s, Got: %s", expResult, format)
	}
}

func TestGetEventsInMinuteError(t *testing.T) {
	expError := fmt.Errorf("failed to parse time string(1214124asfasgasdfgasd)," +
		" layout([2006-01-02 15:04:05]), " +
		"error: parsing time \"1214124asfasgasdfgasd\" " +
		"as \"[2006-01-02 15:04:05]\": " +
		"cannot parse \"1214124asfasgasdfgasd\" as \"[\"")

	tInLayout := "[2006-01-02 15:04:05]"
	tSuffix := "NOK"

	brokenFile, err := os.Open(tBrokenFilePath)
	if err != nil {
		t.Fatalf("Failed to open file path(%s): %s", tBrokenFilePath, err)
	}
	defer brokenFile.Close()

	result, err := getEventsInMinute(
		bufio.NewScanner(brokenFile),
		tInLayout,
		tSuffix,
	)

	if result != nil {
		t.Fatalf("Expected nil results in test, got: %v", result)
	}

	if err.Error() != expError.Error() {
		t.Fatalf("Expected error:%s, Got: %s", expError, err)
	}

}

func TestGetEventsInMinuteResult(t *testing.T) {
	tInLayout := "[2006-01-02 15:04:05]"
	tSuffix := "PIU"

	brokenFile, err := os.Open(tFilePath)
	if err != nil {
		t.Fatalf("Failed to open file path(%s): %s", tFilePath, err)
	}
	defer brokenFile.Close()

	result, err := getEventsInMinute(
		bufio.NewScanner(brokenFile),
		tInLayout,
		tSuffix,
	)
	if err != nil {
		t.Fatalf("Unexpected error in test: %s", err)
	}

	expResult := map[string]int{
		"2018-04-11 03:13:00 +0000 UTC": 1,
		"2018-04-11 03:14:00 +0000 UTC": 0,
		"2018-04-11 03:15:00 +0000 UTC": 2,
		"2018-04-11 04:15:00 +0000 UTC": 1,
		"2018-04-11 04:16:00 +0000 UTC": 0,
		"2018-04-11 04:27:00 +0000 UTC": 1,
	}

	for _, event := range result[1:] {
		key := event.minuteTime.String()
		v, ok := expResult[key]

		if !ok {
			t.Fatalf("no such key in expResult map: %s", key)
		}

		if v != event.count {
			t.Fatalf("Expected %d count for %s time, Got: %d", v, key, event.count)
		}

		delete(expResult, event.minuteTime.String())
	}

	if len(expResult) != 0 {
		t.Fatalf("expResult map is not empty, something gone wrong")
	}

}
