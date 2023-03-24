package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
)

// from https://thedevelopercafe.com/articles/server-sent-events-in-go-595ae2740c7a

//go:embed index.html
var indexHTML []byte

func main() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		ch := make(chan int)
		go generator(r.Context(), ch)
		for data := range ch {
			event, err := formatter("data", data)
			if err != nil {
				log.Println(err)
				break
			}
			_, err = fmt.Fprint(w, event)
			if err != nil {
				log.Println(err)
				break
			}
			flusher.Flush()
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(indexHTML)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
