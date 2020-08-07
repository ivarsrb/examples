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
)

// Define common colors for convenience
var (
	blue  = color.RGBA{0, 0, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
)

func main() {
	// Static files (such as html file) are served from /static folder
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// Handle retrieval of a single jpeg image
	http.HandleFunc("/picture", getPicture)
	// Start a server on the port
	port := "8080"
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getPicture sends jpeg image bytes over http as well as content description
// for browser to able to render the image properly
func getPicture(w http.ResponseWriter, r *http.Request) {
	imgBytes := getJPEG(200, 200, blue)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBytes)))
	w.Write(imgBytes)
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
