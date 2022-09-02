package main

import (
	"fmt"
	"myPhotos/config"
	"myPhotos/entity"
	"myPhotos/third_party/exiftool"
	"myPhotos/web"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("USAGE: myPhotos [CONFIG AND DATA DIRECTORY]")
		os.Exit(1)
	}

	if l, err := os.Lstat(os.Args[1]); err != nil || !l.IsDir() {
		fmt.Println("the dir is not exist or is not a directory")
		os.Exit(1)
	}

	if _, err := exec.LookPath("exiftool"); err != nil {
		fmt.Println("exiftool not found, you can download it from https://exiftool.org")
		os.Exit(1)
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		fmt.Println("ffmpeg not found, you can download it from https://ffmpeg.org")
		os.Exit(1)
	}

	config.InitializeConfig()
	entity.InitializeDatabase()
	exiftool.InitializeExiftool()

	web.Serve()
}
