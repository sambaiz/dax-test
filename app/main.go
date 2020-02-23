package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hi this is ECS")
	})
	log.Printf("Listen on :%s\n", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}
