package main

import (
	"fmt"
	"myPhotos/config"
	"myPhotos/scanner"
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
	}
}

func printUsage() {
	fmt.Println("USAGE:")
	fmt.Println("  myPhotos server              start web server")
	fmt.Println("  myPhotos scan DIR [-g]         scan dir and save to database")
	fmt.Println("                               -g: get geography info(use AMap api)")
}
