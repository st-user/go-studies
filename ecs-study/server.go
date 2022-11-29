package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	port := 80
	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	fmt.Printf("Listening on %d\n", port)
	log.Fatal(server.ListenAndServe())
}
