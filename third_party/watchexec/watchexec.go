package watchexec

import (
	"bufio"
	"fmt"
	"myPhotos/config"
	"myPhotos/entity"
	"myPhotos/logger"
	"myPhotos/services"
	"myPhotos/third_party/exiftool"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WatchExec TODO 目前只执行新增和删除动作 扩展名只支持小写
func WatchExec() {
	path := os.Getenv("WATCHEXEC_COMMON_PATH")

	createdFiles := filepath.SplitList(os.Getenv("WATCHEXEC_CREATED_PATH"))
	removedFiles := filepath.SplitList(os.Getenv("WATCHEXEC_REMOVED_PATH"))
	logger.InfoLogger.Println("created", len(createdFiles), "files")
	logger.InfoLogger.Println("removed", len(removedFiles), "files")

	et, err := exiftool.NewExiftool()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
	defer et.Close()

	for _, file := range createdFiles {
		f := filepath.Join(path, file)
		f = strings.ReplaceAll(f, "\\", "/")
		fm := et.ExtractMetadata(f)
		services.SaveMedia(fm)
	}

	for _, file := range removedFiles {
		f := filepath.Join(path, file)
		f = strings.ReplaceAll(f, "\\", "/")
		entity.RemoveMedia(f)
	}
}

// StartWatchExec TODO 获取 cmd stdout 和 stderr 的部分需要优化
func StartWatchExec(path string) {
	name := "watchexec"
	arg := make([]string, 0)

	dirs := strings.Split(""+path+"", ";")
	dirs = strings.Split(strings.Join(dirs, " -w "), " ")

	extList := append(config.PhotoExtList, config.VideoExtList...)
	arg = append(arg, "--exts", strings.Join(extList, ","), "-w")
	for _, dir := range dirs {
		arg = append(arg, dir)
	}

	path, err := os.Executable()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
	arg = append(arg, path+" watchexec", "--force-poll", "1000")

	cmd := exec.Command(name, arg...)

	// stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
	errScanner := bufio.NewScanner(stderr)
	errScanner.Split(bufio.ScanLines)
	go func() {
		for errScanner.Scan() {
			fmt.Println(errScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
}
