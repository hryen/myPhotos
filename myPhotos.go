package main

import (
	"fmt"
	"myPhotos/config"
	"myPhotos/scanner"
	"myPhotos/third_party/watchexec"
	"myPhotos/web"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}
	if os.Args[1] == "server" {
		web.StartServer()
	} else if os.Args[1] == "scan" {
		config.GPSToGeo = false
		if len(os.Args) == 4 && os.Args[3] == "-g" && config.AMapKey != "" {
			config.GPSToGeo = true
		}
		scanner.StartScan(os.Args[2])
	} else if os.Args[1] == "watchexec" {
		watchexec.WatchExec()
	} else if os.Args[1] == "watch" {
		watchexec.StartWatchExec(os.Args[2])
	}
}

func printUsage() {
	fmt.Println("USAGE:")
	fmt.Println("  homePhoto server              start web server")
	fmt.Println("  homePhoto scan DIR -g         scan dir and save to database")
	fmt.Println("                                -g: get geography info(use AMap api)")
	fmt.Println("  homePhoto watch DIRS...       watch dirs file change, dirs use ; separate")
}
