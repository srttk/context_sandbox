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
	child, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("looping...")

		select {
		case <-child.Done():
			fmt.Println("break loop.")
			return
		default:
			continue
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// WithDeadline 用の time.Timeを作成
	// 現在時刻の3秒後を deadline に設定
	deadline := time.Now().Add(3 * time.Second)
	// 親 context 作成
	parent, cancel := context.WithDeadline(context.Background(), deadline)
	// 上記設定した時間が来たら cancel() を呼ぶ
	defer cancel()

	// ゾンビ goroutine に ctx を渡す
	go leak(parent)

	fmt.Fprint(w, "allow request.")

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
