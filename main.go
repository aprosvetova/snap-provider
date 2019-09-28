package main

import (
	"bytes"
	"fmt"
	"github.com/aprosvetova/snap-provider/mjpeg"
	"github.com/aprosvetova/snap-provider/queue"
	"github.com/cheggaaa/pb/v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var frameQueue *queue.Queue
var bar *pb.ProgressBar

func keepFrames() {
	fmt.Println("Initializing...")
	dec, err := mjpeg.NewDecoderFromURL(mjpegStream)
	if err != nil {
		fmt.Println("Can't open stream", err)
		os.Exit(1)
	}
	systemd := os.Getppid() == 1
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
		if frameQueue.GetLength() < frameBufferLength {
			if systemd {
				updatePlainProgress(frameQueue.GetLength())
			} else {
				updateProgress(frameQueue.GetLength())
			}
		}
		if frameQueue.GetLength() == frameBufferLength-1 {
			fmt.Println("The buffer is full. Ready to use.")
		}
	}
}

func updateProgress(loaded int) {
	if loaded == 1 {
		bar = pb.StartNew(frameBufferLength)
	}
	bar.SetCurrent(int64(loaded) + 1)
	if loaded == frameBufferLength-1 {
		bar.Finish()
	}
}

func updatePlainProgress(loaded int) {
	loaded++
	if (loaded%int(frameBufferLength/10) == 0) || loaded == frameBufferLength {
		fmt.Println("Buffered", loaded, "frames of", frameBufferLength)
	}
}

func saveFrames(w http.ResponseWriter) {
	frames := frameQueue.GetAll()
	cmd := exec.Command("ffmpeg",
		"-framerate", strconv.Itoa(frameRate),
		"-f", "jpeg_pipe",
		"-i", "pipe:0",
		"-c:v", codec,
		"-b:v", bitrate,
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
	parseFlags()

	frameQueue = queue.NewQueue(frameBufferLength)

	go keepFrames()

	http.HandleFunc("/snap.mp4", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "video/mp4")
		t := time.Now()
		saveFrames(w)
		fmt.Println("Given a snap in", time.Since(t))
	})
	http.ListenAndServe(bindAddress, nil)
}
