# Streaming server-side animation to browser using M-JPEG and Go

## Introduction
Most of the time, when developing client-server applications, real-time animations in clients browser are performed
directly on the client side using rendering tools such as WebGL or WebGPU. The data to be rendered and rendering code
are sent from server to client and then rendered by the browser engine using Graphics Processing Unit (GPU) from the client machine. The benefits are obvious - server is freed from potentially massive amount of work if large numbers of clients connect simultaniously. But there are times and places when animating and rendering on the server and sending already rendered images to the client in real-time is exactly what we are after. This tutorial is an introduction to one of the possible solutions to this problem. We are going to use popular Motion JPEG format to transfer data from server to client. Eeven though the idea presented here is language agnostic the example code is written in Golang.  

* What is mjpeg and how to send it over HTTP
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

