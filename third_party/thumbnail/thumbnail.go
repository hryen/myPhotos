package thumbnail

import (
	"myPhotos/config"
	"myPhotos/entity"
	"myPhotos/logger"
	"myPhotos/tools"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func SaveMediaThumbnail(m *entity.Media) {
	thumbnailFile := filepath.Join(config.ThumbnailPath, m.ID+".jpg")
	if _, err := os.Stat(thumbnailFile); !os.IsNotExist(err) {
		return
	}

	ext := strings.ToLower(filepath.Ext(m.Path))

	if tools.ArrayContains(config.PhotoExtList, ext) {
		command := []string{"-i", m.Path, "-vf", "scale=-1:256", thumbnailFile}
		doSaveMediaThumbnail(command)
	}

	if tools.ArrayContains(config.VideoExtList, ext) {
		command := []string{"-ss", "00:00:01.000", "-i", m.Path, "-vframes:v", "1", "-vf", "scale=-1:256", thumbnailFile}
		doSaveMediaThumbnail(command)
	}
}

func doSaveMediaThumbnail(arg []string) {
	cmd := exec.Command("ffmpeg", arg...)
	go func() {
		err := cmd.Run()
		if err != nil {
			logger.ErrorLogger.Println("failed to save thumbnail:", arg, err)
		}
	}()
}
