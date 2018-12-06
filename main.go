package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// 終わらない処理
func leak(ctx context.Context) {
	child, _ := context.WithCancel(ctx)
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("looping...")

		select {
		case <-child.Done():
			fmt.Println("break loop.")
			return
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// 親 context 作成
	parent, cancel := context.WithCancel(context.Background())

	// ゾンビ goroutine に ctx を渡す
	go leak(parent)

	fmt.Fprintf(w, "allow request.")

	// リクエスト返したらゾンビを殺す
	cancel()

	select {
	case <-parent.Done():
		fmt.Println(parent.Err())
	}
}

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
