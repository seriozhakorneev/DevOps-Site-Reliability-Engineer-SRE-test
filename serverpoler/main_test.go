package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPollerError(t *testing.T) {
	t.Parallel()

	testServer := "invalid_server"
	testServers := []string{testServer}
	testResult := make(chan string)

	expectedOut := fmt.Sprintf(
		"invalid_server request failed: "+
			"Get \"%s\": "+
			"dial tcp: lookup invalid_server: no such host",
		httpPrefix+testServer+metricPath,
	)

	go poller(testServers, time.Second*1, testResult)

	out := strings.TrimSpace(<-testResult)

	if out[20:] != expectedOut {
		t.Fatalf("expected output: %s, got: %s", expectedOut, out)
	}
}

func TestPollerResult(t *testing.T) {
	t.Parallel()

	expCount := 220

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf(`{"count": %d}`, expCount)))
	}))
	defer server.Close()

	testResult := make(chan string)
	testServer := server.URL[7:]

	go poller([]string{testServer}, time.Second*1, testResult)

	out := <-testResult
	expectedOut := fmt.Sprintf("%s %d", testServer, expCount)

	if out[20:] == expectedOut {
		t.Fatalf("expected output: %s, got: %s", expectedOut, out)
	}

}

func TestGetCountRequestFailed(t *testing.T) {
	t.Parallel()

	expResult := 0
	expError := errors.New("request failed: Get \"invalid_url\": unsupported protocol scheme \"\"")

	count, err := getCount("invalid_url")
	if count != expResult {
		t.Fatalf("Expected count: %d, got: %d", expResult, count)
	}

	if err.Error() != expError.Error() {
		t.Fatalf("Expected err: %s\ngot: %s", expError, err)
	}
}

func TestGetCountStatusNotOK(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(``))
	}))
	defer server.Close()

	expResult := 0
	expError := errors.New("response status code is not 200: Status Code: 400")

	count, err := getCount(server.URL)
	if count != expResult {
		t.Fatalf("Expected count: %d, got: %d", expResult, count)
	}

	if !reflect.DeepEqual(err, expError) {
		t.Fatalf("Expected err: %s\ngot: %s", expError, err)

	}
}

func TestGetCountStatusNotJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
		rw.Write([]byte(``))
	}))
	defer server.Close()

	expResult := 0
	expError := errors.New("content-type header is not " +
		"application/json: Content-Type: text/html;charset=UTF-8")

	count, err := getCount(server.URL)
	if count != expResult {
		t.Fatalf("Expected count: %d, got: %d", expResult, count)
	}

	if !reflect.DeepEqual(err, expError) {
		t.Fatalf("Expected err: %s\ngot: %s", expError, err)

	}
}

func TestGetCountStatusDecodeJSONFailed(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(``))
	}))
	defer server.Close()

	expResult := 0
	expError := errors.New("decode json failed: EOF")

	count, err := getCount(server.URL)
	if count != expResult {
		t.Fatalf("Expected count: %d, got: %d", expResult, count)
	}

	if !reflect.DeepEqual(err.Error(), expError.Error()) {
		t.Fatalf("Expected err: %s\ngot: %s", expError, err)
	}
}

func TestGetCountStatusDecodeCountIsNil(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{}`))
	}))
	defer server.Close()

	expResult := 0
	expError := fmt.Errorf("response data is empty")

	count, err := getCount(server.URL)
	if count != expResult {
		t.Fatalf("Expected count: %d, got: %d", expResult, count)
	}

	if !reflect.DeepEqual(err.Error(), expError.Error()) {
		t.Fatalf("Expected err: %s\ngot: %s", expError, err)
	}
}

func TestGetCountResults(t *testing.T) {
	t.Parallel()

	for i := 0; i <= 1000; i++ {
		jsonStr := fmt.Sprintf(`{"count": %d}`, i)

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(jsonStr))
		}))

		expResult := i
		count, err := getCount(server.URL)
		if err != nil {
			t.Fatalf("Unexpected error in testing: %s", err)
		}

		if count != expResult {
			t.Fatalf("Expected count: %d, got: %d", expResult, count)
		}

		server.Close()
	}
}
