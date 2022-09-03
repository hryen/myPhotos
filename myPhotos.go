package main

import (
	"fmt"
	"myPhotos/config"
	"myPhotos/entity"
	"myPhotos/logger"
	"myPhotos/third_party/exiftool"
	"myPhotos/web"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("USAGE: myPhotos [CONFIG AND DATA DIRECTORY]")
		os.Exit(1)
	}

	if l, err := os.Lstat(os.Args[1]); err != nil || !l.IsDir() {
		fmt.Println("the directory is not exist or is not a directory")
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

	// 优雅的结束程序
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
	}()

	// 启动程序
	go func() {
		defer wg.Done()

		config.InitializeConfig(os.Args[1])
		entity.InitializeDatabase()
		exiftool.InitializeExiftool()

		web.StartServer()
	}()

	wg.Wait()

	terminate()
}

func terminate() {
	logger.InfoLogger.Println("terminating...")

	err := entity.Close()
	if err != nil {
		logger.ErrorLogger.Println("close database error:", err)
	}
	logger.InfoLogger.Println("closed database")

	// TODO
	// close exiftool error: error while closing exiftool: [error while waiting for exiftool to exit: exit status 0xc000013a]
	err = exiftool.Et.Close()
	if err != nil {
		logger.ErrorLogger.Println("close exiftool error:", err)
	}

}
