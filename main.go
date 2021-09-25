package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

const dst = "https://en867p8z3qcic.x.pipedream.net/"

func handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing required \"name\" query parameter", http.StatusBadRequest)
			return
		}

		body := fmt.Sprintf("name=%s", name)
		//body := url.Values{"name": {name}}.Encode()

		resp, err := http.Post(dst, "application/x-www-form-urlencoded", bytes.NewReader([]byte(body)))
		if err != nil {
			log.Printf("failed to post: %v", err)
			http.Error(w, "failed to post", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			log.Printf("failed to discard body: %v", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	http.Handle("/", handle())
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))
}
