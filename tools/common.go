package tools

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func ArrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetModTime(path string) time.Time {
	lstat, _ := os.Lstat(path)
	return lstat.ModTime()
}

func NumberToString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch ret := v.(type) {
	case string:
		return ret
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", ret)
	default:
		return ""
	}
}

func ConvertDuration(duration float64) string {
	d := fmt.Sprintf("%v", duration)

	seconds, _ := strconv.Atoi(strings.Split(d, ".")[0])
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = seconds % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = seconds % 60

	t := ""
	if hours > 0 {
		t += fmt.Sprintf("%.0f", hours) + ":"
	}
	if minutes > 0 {
		t += fmt.Sprintf("%.0f", minutes) + ":"
	} else {
		t += "0:"
	}
	if seconds < 10 {
		t += fmt.Sprintf("0%d", seconds)
	} else {
		t += fmt.Sprintf("%d", seconds)
	}

	return t
}
