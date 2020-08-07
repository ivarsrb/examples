package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/render", render)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func render(w http.ResponseWriter, r *http.Request) {

}
