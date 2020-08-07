# Streaming server-side animation to browser using M-JPEG and Go

## Introduction
Most of the time, when developing client-server applications, real-time animations are performed directly on the client side using rendering APIs such as WebGL or WebGPU. The data to be rendered and rendering code are sent from server to client and then rendered by the browser engine using Graphics Processing Unit (GPU) from the clients machine. The benefits are obvious - server is freed from potentially massive amount of work if large numbers of clients connect simultaniously.  
But there's a time and place when rendering on the server and streaming already rendered images to the client in real-time is exactly what we are after. This tutorial is an introduction to one of the possible solutions to this problem.  
We are going to use popular Motion-JPEG format to stream an animation from server to client. On our server we will procedurally generate sequence of images that will correspond to frames in our animation. Even though the technique presented here is language agnostic the example code is written in Golang. Source code for this tutorial can be found at ...
## What is Motion JPEG?
Motion JPEG or M-JPEG is a video format in which each frame is a single JPEG image. It is widely used format in video-capturing devices and most major web browsers also support it. It may not be the most efficient video format but it has an advantage of being easy to use and understand. M-JPEG is selected for this tutorial exactly because of that.
[iamge of separate jpeg frames near each other]
## Sending M-JPEG over HTTP
Essentially what we are going to do is stream image-by-image from server to client over the widely used HTTP (Hypertext Transfer Protocol) protocol.   
In the browser images will be displayed in `<img>` HTML element one after the other giving an illusion of animation. Images will be recieved from the *URL* provided in *src* attribute. Since we are going to generate content on-the-fly the *URL* will be set to some path to which *GET* request is going to be sent by the browser. Server is going to respond to this request by sending back image data as well as appropriate content headers.  
When sending M-JPEG over HTTP the servers response should include:  

    Content-Type: multipart/x-mixed-replace; boundary="<boundary-name>"

    --<boundary-name>
    Content-Type: image/jpeg
    Content-Length: <number of bytes>

    <JPEG image_1 bytes>
    --<boundary-name>
    Content-Type: image/jpeg
    Content-Length: <number of bytes>

    <JPEG image_2 bytes>
    --<boundary-name>
    Content-Type: image/jpeg
    Content-Length: <number of bytes>

    <JPEG image_n bytes>
    --<boundary-name>

`Content-Type: multipart/x-mixed-replace; boundary="<boundary-name>"` tells our client that response is going to come in multiple parts (images) that are separated by `<boundary-name>` and will be replaced each time.  
`Content-Type: image/jpeg` indicates that browser should interpret recieved byte data as JPEG image and `Content-Length: <number of bytes>` should be the size of this data.  
`--<boundary-name>` points to where each image (frame in our case) starts and ends and `<JPEG image_n bytes>` of course is the image data istelf.  
Note that blank lines between some parts of the response are important.  

## Creating JPEG images in Go
From the HTTP response format above we can observe that actual JPEG image data is the *key piece in the puzzle*, so, as a first step, let's create a JPEG image in memory and fill it with color data.  
How do we do that? Luckally for us, Go provides us with *Image* package from it's standard library which we will use now.  
We will define our sample image using `Image.RGBA` type from *Image* package; Create image of the given size
``` go
im := image.NewRGBA(image.Rect(0, 0, w, h))
```
and fill it with a color.
``` go
draw.Draw(im, im.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
```
The color is described with `color.RGBA` struct type.  
At this moment the image is not of JPEG format also note that what we are interested in right now is a JPEG image data represented as a buffer of bytes since we want to send it over HTTP to client browser. Let's do exactly that - encode our image as JPEG byte buffer
``` go
var buff bytes.Buffer
jpeg.Encode(&buff, im, nil)
```
We might want to wrap this code in convinience function that takes iamge parameters, like dimensions and color, and returns pointer to byte data.
``` go
func getJPEG(w int, h int, color color.RGBA) []byte {
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(im, im.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	var buff bytes.Buffer
	jpeg.Encode(&buff, im, nil)
	return buff.Bytes()
}
```
And use it like
``` go
imgBytes := getJPEG(200, 200, color.RGBA{0, 0, 255, 255})
```
Now that we have image data in memory and pointer pointing to it let's see how to send it over HTTP connection.
## Sending image over HTTP and rendering it in browser
* create single blue image and show it in browser 
## Sending images one after the other
* Create simple animation of changing three images in browser
## Sine wave animation
* Explain sine equation
* Show how to render it live

