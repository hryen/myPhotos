package thumbnail

import (
	"myPhotos/config"
	"myPhotos/logger"
	"myPhotos/tools"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func SaveMediaThumbnail(file string) {
	filename := filepath.Base(file)
	thumbnailFile := filepath.Join(config.ThumbnailPath, filename) + ".thumbnail.jpg"
	if _, err := os.Stat(thumbnailFile); !os.IsNotExist(err) {
		return
	}

	ext := strings.ToLower(filepath.Ext(file))

	if tools.ArrayContains(config.PhotoExtList, ext) {
		command := []string{"-i", file, "-vf", "scale=-1:256", thumbnailFile}
		doSaveMediaThumbnail(command)
	}

	if tools.ArrayContains(config.VideoExtList, ext) {
		command := []string{"-ss", "00:00:01.000", "-i", file, "-vframes:v", "1", "-vf", "scale=-1:256", thumbnailFile}
		doSaveMediaThumbnail(command)
	}
}

func doSaveMediaThumbnail(arg []string) {
	cmd := exec.Command("ffmpeg", arg...)
	err := cmd.Start()
	if err != nil {
		logger.ErrorLogger.Println("failed to save thumbnail:", arg, err)
	}
}
