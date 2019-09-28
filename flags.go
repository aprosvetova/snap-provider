package main

import (
	"flag"
	"os"
)

var mjpegStream string
var frameBufferLength int
var frameRate int
var codec string
var bitrate string
var bindAddress string

func parseFlags() {
	flag.StringVar(&mjpegStream, "stream", "", "full URL to mjpeg stream")
	flag.IntVar(&frameBufferLength, "buffer", 150, "number of frames to be buffered")
	flag.IntVar(&frameRate, "fps", 15, "output video framerate")
	flag.StringVar(&codec, "codec", "libx264", "mp4 ffmpeg codec (-c:v)")
	flag.StringVar(&bitrate, "bitrate", "1000K", "mp4 ffmpeg bitrate (-b:v)")
	flag.StringVar(&bindAddress, "bind", "0.0.0.0:5533", "http server bind address")

	flag.Parse()

	if mjpegStream == "" || frameBufferLength <= 0 || frameRate <= 0 || codec == "" ||
		bitrate == "" || bindAddress == "" {
		flag.Usage()
		os.Exit(2)
	}
}
