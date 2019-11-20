# Snap Provider

This is a simple daemon that listens to MJPEG camera stream, buffers last N frames and serves as a video upon request on `http://<bindAddress>/snap.mp4` endpoint.

Can be configured with environment variables. All vars are required.

|Var|Description|Default value|
|-|-|-|
|STREAM|full URL to mjpeg stream||
|BUFFER|how many frames to buffer|150|
|FPS|frame rate of resulting video|15|
|CODEC|FFmpeg codec to be used for video encoding|libx264|
|BITRATE|resulting video bitrate, passed to -b:v of ffmpeg|1000K|
|BIND|HTTP server bind address|0.0.0.0:5533|

Resulting video will be `frameBufferLength/frameRate` seconds length.
If you have 15 FPS stream and want 10 seconds of video, use `-buffer 150 -fps 15` flags.