package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func heavy(result chan<- string) {
	// 5秒かかる
	fmt.Println("process started.")
	time.Sleep(5 * time.Second)
	fmt.Println("process finished.")
	// 処理が完了したらチャネルに値を送信する
	result <- "process succeeded!"
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 結果を受信するチャネル
	result := make(chan string, 1)
	// 重い処理を起動
	go heavy(result)

	// ...省略

	select {
	case r := <-result:
		fmt.Fprintf(w, "allow request. result: %v\n", r)
	}
}

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
