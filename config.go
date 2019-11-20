package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"os"
)

type config struct {
	MjpegStream       string `env:"STREAM"`
	FrameBufferLength int    `env:"BUFFER" envDefault:"150"`
	FrameRate         int    `env:"FPS" envDefault:"15"`
	Codec             string `env:"CODEC" envDefault:"libx264"`
	Bitrate           string `env:"BITRATE" envDefault:"1000K"`
	BindAddress       string `env:"BIND" envDefault:"0.0.0.0:5533"`
}

func loadConfig() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Println("can't parse env")
		os.Exit(2)
	}
	if cfg.MjpegStream == "" || cfg.FrameBufferLength <= 0 || cfg.FrameRate <= 0 || cfg.Codec == "" ||
		cfg.Bitrate == "" || cfg.BindAddress == "" {
		fmt.Println("wrong env values")
		os.Exit(2)
	}

	return cfg
}
