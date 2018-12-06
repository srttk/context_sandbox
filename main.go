package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// 終わらない処理
func leak(done <-chan struct{}) {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("looping...")

		select {
		case <-done:
			fmt.Println("break loop.")
			return
		default:
			continue
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// プロセス終了を伝えるチャネルを作成
	done := make(chan struct{})
	// done を渡す
	go leak(done)

	fmt.Fprint(w, "allow request.")

	select {
	case <-time.After(3 * time.Second):
		close(done)
	}
}

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
