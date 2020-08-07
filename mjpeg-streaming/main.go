// Package main shows an example for transferring jpeg stream over HTTP
package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Define common colors for convenience
var (
	blue  = color.RGBA{0, 0, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
)

// Boundary will separate frames in M-JPEG animation transfer
const boundary = "abcd4321"

func main() {
	// Static files (such as html file) are served from /static folder
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// Handle retrieval of a single jpeg image
	http.HandleFunc("/picture", getPicture)
	// Handle simple animation request
	http.HandleFunc("/animation", getAnimation)
	// Start a server on the port
	port := "8080"
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getJPEG creates a single color image with given dimensions and color.
// Returns the image as a slice of jpeg bytes
func getJPEG(w int, h int, color color.RGBA) []byte {
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(im, im.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	var buff bytes.Buffer
	jpeg.Encode(&buff, im, nil)
	return buff.Bytes()
}

// getPicture sends jpeg image bytes over http as well as content description
// for browser to able to render the image properly
func getPicture(w http.ResponseWriter, r *http.Request) {
	imgBytes := getJPEG(200, 200, blue)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBytes)))
	w.Write(imgBytes)
}

// getAnimation creates sample images and sends them one after the other to client
func getAnimation(w http.ResponseWriter, r *http.Request) {
	// Sample images
	imgBlue := getJPEG(200, 200, blue)
	imgRed := getJPEG(200, 200, red)
	imgGreen := getJPEG(200, 200, green)
	delay := 500 * time.Millisecond
	// To send buffered data to client
	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("Buffer fulshing is not implemented")
	}
	// Set headers and content to send as a response
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)

	w.Write([]byte("\r\n--" + boundary + "\r\n"))
	w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgBlue)) + "\r\n\r\n"))
	w.Write(imgBlue)
	w.Write([]byte("\r\n--" + boundary + "\r\n"))

	// Otherwise buffer will be flushed after handler exits or buffer maxsize is full
	f.Flush()
	// Delay
	time.Sleep(delay)

	w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgRed)) + "\r\n\r\n"))
	w.Write(imgRed)
	w.Write([]byte("\r\n--" + boundary + "\r\n"))

	f.Flush()
	time.Sleep(delay)

	w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgGreen)) + "\r\n\r\n"))
	w.Write(imgGreen)
	w.Write([]byte("\r\n--" + boundary + "\r\n"))
}
