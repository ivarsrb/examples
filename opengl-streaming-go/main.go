// Package main is an example of OpenGL rendered content streaming over HTTP
package main

import (
	"log"
	"net/http"
)

// Boundary will separate frames in M-JPEG animation transfer
const boundary = "abcd4321"

func main() {
	// Static files (such as html file) are served from /static folder
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// Handle retrieval of a opengl rendered content in the form of MJPEG
	http.HandleFunc("/opengl", getOpenglContent)
	// Start a server on the port
	port := "8080"
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getOpenglContent renders an OpenGL scene and sends the rendered images over to client
// in the form of MJPEG video
func getOpenglContent(w http.ResponseWriter, r *http.Request) {

}
