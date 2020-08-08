# Server-side animation live streaming with MJPEG and Go

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
And use it like this to create blue image of 200x200 pixels
``` go
imgBytes := getJPEG(200, 200, color.RGBA{0, 0, 255, 255})
```
Now that we have an image data in memory and a pointer pointing to it let's look at how to send it over HTTP connection.
## Sending image over HTTP to client browser
In order to create our client-server application we will ned bot server and client parts. Th client par will be simple HTML document, for the server part we will create HTTP server that lsitenss to incomming connections and responds with imaga data when requested.
In `main` function add the following lines
``` go
func main() {
    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.HandleFunc("/picture", getPicture)
    port := "8080"
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
```
HTTP server wil listen for incomming connections on port `8080`. Static files, like HTML index page will be served from `/static` directory. If the client sends GET request to `/picture` URL path the response will be sent back by our server. In this case function `getPicture` is responsible for sending back our image data.  
Client code is basic HTML template with `<img>` element that sends GET request to URL `/pictures` upon page loading and renders the response as an image.
``` HTML
<body>
    <img style="border:2px solid black" src="/picture"  />
</body>
```
Our job is to send the response in the way browser can correctly interpret. Let's implement `getPicture` function.
``` go
func getPicture(w http.ResponseWriter, r *http.Request) {
	imgBytes := getJPEG(200, 200, blue)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBytes)))
	w.Write(imgBytes)
}
```
At first we create image in memory and return pointer to slice of bytes. Then we set our response headers in the way that lets browser understand what ir gets back. We set image type and image size. And finally we send the byte buffer itself.
Spin up the server with command `go run` and open browser pointing to `localhost:8080`. You should see blue square that is effectivly the image we crated and sent from server. "But there is no animation in this example!" you say. Let's fix this now! 
## Simple animation
Let's add new URL path to our server and a handler respensible for the response.
``` go
...
http.HandleFunc("/animation", getAnimation)
...
```
``` go
func getAnimation(w http.ResponseWriter, r *http.Request) {
}
```
Also don't forget to change `src` attribute in `<img>` tag.   
To create basic animation we will need a couple images to show one right after another. By using our helper function create three images - red, yellow and green, that, when animated, will give as an illusion of changing traffic lights.
``` go
size = 200
var (
	red    = color.RGBA{255, 0, 0, 255}
	green  = color.RGBA{0, 255, 0, 255}
	yellow = color.RGBA{255, 255, 0, 255}
)
imgRed := getJPEG(size, size, red)
imgYellow := getJPEG(size, size, yellow)
imgGreen := getJPEG(size, size, green)
```
Now we can start sending animation back to our client. At first we indicate that response is going to consist of multiple parts separated by `boundry`.
``` go
const boundary = "abcd4321"
w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
```
For `boundry` I chose an arbitrary string, but it could be anything as long as it is not going to appear in the data we want to separate.  
At this point, according to M-JPEG response format we discussed earlier, we can stream all our images one by one and hopefully see the animation.
``` go
w.Write([]byte("\r\n--" + boundary + "\r\n"))
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgRed)) + "\r\n\r\n"))
w.Write(imgRed)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgYellow)) + "\r\n\r\n"))
w.Write(imgYellow)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgGreen)) + "\r\n\r\n"))
w.Write(imgGreen)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
```
When you run the server and refresh your browser you will only see a green square. Why? The reason is that all of our response is sent all at once instead of image by image. This is how HTTP originally operates - the response is gathered and sent when it is finished or buffer size limit is reached. But this is not what we want, we are developing a live-streaming service. Luckally for us we can work around that. `http.ResponseWriter` parameter that is provided to us through the handler (usually) implements `http.Flusher` that will allow us to flsuh the buffer and send our data to client immidiatly. Let's obtain `Flusher` like so
``` go
f, ok := w.(http.Flusher)
if !ok {
    log.Println("HTTP buffer flushing is not implemented")
}
```
and call it's `Flush()` method in between frames 
``` go
w.Write([]byte("\r\n--" + boundary + "\r\n"))
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgRed)) + "\r\n\r\n"))
w.Write(imgRed)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
f.Flush()
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgYellow)) + "\r\n\r\n"))
w.Write(imgYellow)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
f.Flush()
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgGreen)) + "\r\n\r\n"))
w.Write(imgGreen)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
```
When you re-run the server now you will still probably see a green square with no animation. There is one more thing missing that you probably already guessed - we have no delay between our frames, so we just physically can't catch the sight of the red and yellow images as they whiz by. Let's insert delay in between our frames
``` go
delay = 500 * time.Millisecond
w.Write([]byte("\r\n--" + boundary + "\r\n"))
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgRed)) + "\r\n\r\n"))
w.Write(imgRed)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
f.Flush()
time.Sleep(delay)
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgYellow)) + "\r\n\r\n"))
w.Write(imgYellow)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
f.Flush()
time.Sleep(delay)
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgGreen)) + "\r\n\r\n"))
w.Write(imgGreen)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
```
Finally you should be able to see an animation we were aiming at.  
It worked fine as a demonstration but obviously *hardcoding* frame after frame is not the proper way to write and animation. In the next section we will implement a procedural approach that will be a lot more flexible. 
## Sine wave animation
In this section we are going to make something a bit more interesting than traffic light. Let's animate a sliding wave. Wave will be simulated using simple sine wave function:
> **y** = A * sin(B***x** + C) + D 

where `A` is an amplitude, `B` is period, `C` is phase shift and `D` is vertical shift. `x` and `y` are coordinates in our image space. Since we want to animate the wave we will change phase shift `C` each frame so the wave appears to be sliding right to left. 
[image of sine wave function]  
In terms of **pseudocode** our animation algorithm can be described as follows
```
for number_of_frames as t:
    img = new_image
    for width_of_image as n:
        x = n
        y = A*sin(B*x + t) + D
        img.set_pixel( x, y, color )
    send_to_client( img )
```
Having defined totol number of frames in animation, we start each frame with a blank image, then for all horizontal image pixel we find corresponding vertical pixel and render it in some other color, then, when image is formed, we send it over HTTP to client and repeat.  
Go ahead and create new URL path `/wave` and add new handler, changle `src` atribute in HTML document to `="/wave"`
``` go
...
http.HandleFunc("/wave", getSinewaves)
...
```
``` go
func getSinewaves(w http.ResponseWriter, r *http.Request) {
}
```
The implementation is straight-forward. The only real difference here is that we send images over HTTP in a loop instead of hard-coding each frame in code.  
At first we define parameters that will affect our animation
``` go
const (
    width  = 400
    height = 300
    nframes = 60
    delay = 50 * time.Millisecond
)
```
we are familiar with the dimensions and delay; `nframes` is a number of frames in our animation.
``` go
w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
for t := 0; t < nframes; t++ {
 ...
}
```
Just like in previous example set content type header to indicate that we are about to send a stream in multiple parts. After that start out animation loop. Inside the loop create a new image like so
``` go
img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
```
this is a new image that consists of color pointed to in `palette` argument, so we need to also define our palette
``` go
var palette = []color.Color{color.White, blue}
// Indexes in palette
const (
    whiteIndex = 0
    blueIndex  = 1
)
```
After we have our image we can start rendering to it. Start another loop that loops across horizontal image dimension
``` go
for n := 0; n < width; n++ {
    x := float64(n)
    a := height / 3.0
    b := 0.01
    c := float64(t) / 6.0
    d := height / 2.0
    y := a*math.Sin(x*b+c) + d
    img.SetColorIndex(int(x), int(y), blueIndex)
}
```
At first we calculate our `x` and `y` coordinates that are to be colored according to sine function discussed earlier and then we set calculated pixel in the image `img` at the coordinates `(x, y)` with a color at index `blueIndex` from our image color palette.  
The rest of the code is similar to one we discussed in the previous section. At first we encode an image into jpeg format and retrieve a slice of it's byte array;
``` go
var buff bytes.Buffer
jpeg.Encode(&buff, img, nil)
imgBytes := buff.Bytes()
```
then we form our response contents. Remember that animation should start with boundary marker, otherwise our first frame wil be skipped;
``` go
if t == 0 {
    w.Write([]byte("\r\n--" + boundary + "\r\n"))
}
w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgBytes)) + "\r\n\r\n"))
w.Write(imgBytes)
w.Write([]byte("\r\n--" + boundary + "\r\n"))
``` 
and finally we send our buffered image contents to client as well as pausing a bit bat the end of each frame
``` go
f.Flush()
time.Sleep(delay)
```
Overall the handler function should look like
``` go
func getSinewaves(w http.ResponseWriter, r *http.Request) {
	var palette = []color.Color{color.White, blue}
	const (
		whiteIndex = 0
		blueIndex  = 1
	)
	const (
		width  = 400
		height = 300
		nframes = 60
		delay = 50 * time.Millisecond
	)
	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("HTTP buffer flushing is not implemented")
	}
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
	for t := 0; t < nframes; t++ {
		img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
		for n := 0; n < width; n++ {
			x := float64(n)
			a := height / 3.0
			b := 0.01
			c := float64(t) / 6.0
			d := height / 2.0
			y := a*math.Sin(x*b+c) + d
			img.SetColorIndex(int(x), int(y), blueIndex)
		}
		var buff bytes.Buffer
		jpeg.Encode(&buff, img, nil)
		imgBytes := buff.Bytes()
		if t == 0 {
			w.Write([]byte("\r\n--" + boundary + "\r\n"))
		}
		w.Write([]byte("Content-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(imgBytes)) + "\r\n\r\n"))
		w.Write(imgBytes)
		w.Write([]byte("\r\n--" + boundary + "\r\n"))
		f.Flush()
		time.Sleep(delay)
	}
}
```
Now when you start the server and refresh your browser you should see sine wave animation.
## Conclusion
We've created a simple animation by generating JPEG images on the server and streaming them one by one to the client which in turn rendered them to browser.  
The code above is for presenting the live-streaming idea and is in no way optimised. Go standard library `Image` package probably isn't well suitet for real time animation either. In the next tutorial we will try streaming something more colorfull using a more powerfull rendering library - OpenGL.
