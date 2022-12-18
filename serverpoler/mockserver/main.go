package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func count(w http.ResponseWriter, _ *http.Request) {
	jsonResponse, _ := json.Marshal(map[string]int{"count": 42})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/api/count", count)

	port := ":2020"
	log.Println("listen and serve on:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
