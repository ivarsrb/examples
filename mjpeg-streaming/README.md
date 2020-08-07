# Streaming server-side animation to browser using M-JPEG and Go

## Introduction
Most of the time, when developing client-server applications, real-time animations are performed directly on the client side using rendering APIs such as WebGL or WebGPU. The data to be rendered and rendering code are sent from server to client and then rendered by the browser engine using Graphics Processing Unit (GPU) from the clients machine. The benefits are obvious - server is freed from potentially massive amount of work if large numbers of clients connect simultaniously.  
But there are times and places when animating and rendering on the server and streaming already rendered images to the client in real-time is exactly what we are after. This tutorial is an introduction to one of the possible solutions to this problem.  
We are going to use popular Motion-JPEG format to transfer data from server to client. Even though the technique presented here is language agnostic the example code is written in Golang.  
## What is Motion JPEG?
Motion JPEG or M-JPEG is a video format in which each frame is a single JPEG image. It is widely used format in video-capturing devices and most major web browsers also support it. It may not be the most efficient video format but it has an advantage of being easy to use and understand. M-JPEG is selected for this tutorial exactly because of that.
[iamge of separate jpeg frames near each other]
## Sending M-JPEG over HTTP
Essentially what we are going to do is stream image-by-image from server to client over the widely used HTTP (Hypertext Transfer Protocol) protocol.   
In a browser images will be displayed in `<img>` HTML element one after the other giving an illusion of animation. Images will be recieved from the *URL* provided in *src* attribute. Since we are going to generate content on-the-fly the *URL* will be set to some path to which *GET* request is going to be sent by the browser. Server is going to respond to this request by sending back image data as well as appropriate content headers.  
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
`Content-Type: image/jpeg` tells that browser should expect JPEG data and `Content-Length: <number of bytes>` will be the size of this data.  
`--<boundary-name>` points to where each image (frame in our case) starts and ends and `<JPEG image_n bytes>` of course is the image data istelf.  
Note that blank lines between some parts of the response are important. 

## Creating jpeg image in Go
* How to create jpeg image on the fly and fill it with color. Probably references will be from golang image tutorial.
* Code snippets.
## Sending image over HTTP and rendering it in browser
* create single blue image and show it in browser 
## Sending images one after the other
* Create simple animation of changing three images in browser
## Sine wave animation
* Explain sine equation
* Show how to render it live

