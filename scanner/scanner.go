package scanner

import (
	"io/fs"
	"myPhotos/config"
	"myPhotos/logger"
	"myPhotos/services"
	"myPhotos/third_party/exiftool"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func StartScan(dir string) {
	defer timeCost(time.Now())

	files, err := ScanMediaFile(dir)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	et, err := exiftool.NewExiftool()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
	defer et.Close()

	for _, file := range files {
		fm := et.ExtractMetadata(file)
		services.SaveMedia(fm)
	}

	logger.InfoLogger.Println("Scanned", len(files), "files")
}

func ScanMediaFile(dir string) ([]string, error) {
	files := make([]string, 0)
	extList := append(config.PhotoExtList, config.VideoExtList...)
	err := filepath.WalkDir(dir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !de.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			flag := false
			for _, e := range extList {
				if ext == e {
					flag = true
					break
				}
			}
			if flag {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

func timeCost(start time.Time) {
	tc := time.Since(start)
	logger.InfoLogger.Printf("Time consuming: %v\n", tc)
}
