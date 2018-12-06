package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// 終わらない処理
func leak() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("looping...")
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// ゾンビ goroutine
	go leak()

	fmt.Fprintf(w, "allow request.")
}

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
