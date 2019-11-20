package main

import (
	"bytes"
	"fmt"
	"github.com/aprosvetova/snap-provider/mjpeg"
	"github.com/aprosvetova/snap-provider/queue"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var frameQueue *queue.Queue
var conf config

func keepFrames() {
	fmt.Println("Initializing...")
	dec, err := mjpeg.NewDecoderFromURL(conf.MjpegStream)
	if err != nil {
		fmt.Println("Can't open stream", err)
		os.Exit(1)
	}
	for {
		p, err := dec.GetPart()
		if err != nil {
			fmt.Println("Can't get part", err)
			continue
		}
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(p)
		if err != nil {
			fmt.Println("Can't read part", err)
			continue
		}
		frameQueue.Push(buf.Bytes())
		if frameQueue.GetLength() < conf.FrameBufferLength {
			updateProgress(frameQueue.GetLength())
		}
		if frameQueue.GetLength() == conf.FrameBufferLength-1 {
			fmt.Println("The buffer is full. Ready to use.")
		}
	}
}

func updateProgress(loaded int) {
	loaded++
	if (loaded%int(conf.FrameBufferLength/10) == 0) || loaded == conf.FrameBufferLength {
		fmt.Println("Buffered", loaded, "frames of", conf.FrameBufferLength)
	}
}

func saveFrames(w http.ResponseWriter) {
	frames := frameQueue.GetAll()
	cmd := exec.Command("ffmpeg",
		"-framerate", strconv.Itoa(conf.FrameRate),
		"-f", "jpeg_pipe",
		"-i", "pipe:0",
		"-c:v", conf.Codec,
		"-b:v", conf.Bitrate,
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4",
		"-an",
		"pipe:1")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	go io.Copy(w, stdout)
	for _, frame := range frames {
		io.Copy(stdin, bytes.NewReader(frame))
	}
	stdin.Close()
	cmd.Wait()
}

func main() {
	conf = loadConfig()

	frameQueue = queue.NewQueue(conf.FrameBufferLength)

	go keepFrames()

	http.HandleFunc("/snap.mp4", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "video/mp4")
		t := time.Now()
		saveFrames(w)
		fmt.Println("Given a snap in", time.Since(t))
	})
	http.ListenAndServe(conf.BindAddress, nil)
}
