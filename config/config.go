package config

import (
	"github.com/spf13/viper"
	"myPhotos/logger"
	"os"
	"path/filepath"
	"time"
)

// DataPath 数据存放路径，包含缩略图、数据库等文件
var DataPath string

// ThumbnailPath 保存照片和视频缩略图的路径
var ThumbnailPath string

// UploadPath 保存上传的照片和视频的路径
var UploadPath string

// GPSToGeo 是否开启通过 GPS 经纬度数据使用高德地图 api 反查地理位置
var GPSToGeo bool

// AMapKey 高德地图 key
var AMapKey string

// PhotoExtList 照片的扩展名列表，小写
var PhotoExtList []string

// VideoExtList 视频的扩展名列表，小写
var VideoExtList []string

// Location 时区
var Location *time.Location

// HTTPAddr web服务监听地址
var HTTPAddr string

// HTTPPort web服务监听端口
var HTTPPort string

func InitializeConfig(path string) {
	DataPath = path

	// 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(DataPath)
	err := viper.ReadInConfig()
	if err != nil {
		logger.ErrorLogger.Println("fatal error config file:", err)
		os.Exit(1)
	}

	ThumbnailPath = filepath.Join(DataPath, "thumbnails")
	UploadPath = filepath.Join(DataPath, "uploads")
	GPSToGeo = viper.GetBool("GPSToGeo")
	AMapKey = viper.GetString("AMapKey")
	if AMapKey == "" {
		GPSToGeo = false
	}
	PhotoExtList = viper.GetStringSlice("PhotoExtList")
	VideoExtList = viper.GetStringSlice("VideoExtList")
	HTTPAddr = viper.GetString("HTTPAddr")
	HTTPPort = viper.GetString("HTTPPort")

	// 创建缩略图目录
	if s, err := os.Stat(ThumbnailPath); os.IsNotExist(err) {
		err := os.MkdirAll(ThumbnailPath, os.ModeDir)
		if err != nil {
			logger.ErrorLogger.Println("failed to create thumbnail directory", err)
			os.Exit(1)
		}
	} else {
		if !s.IsDir() {
			logger.ErrorLogger.Println("failed to create thumbnail directory, file is exist but is not directory.")
			os.Exit(1)
		}
	}

	// 创建上传目录
	if s, err := os.Stat(UploadPath); os.IsNotExist(err) {
		err := os.MkdirAll(UploadPath, os.ModeDir)
		if err != nil {
			logger.ErrorLogger.Println("failed to create upload directory", err)
			os.Exit(1)
		}
	} else {
		if !s.IsDir() {
			logger.ErrorLogger.Println("failed to create upload directory, file is exist but is not directory.")
			os.Exit(1)
		}
	}

	// 时区
	Location, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.ErrorLogger.Println("failed init time location", err)
		os.Exit(1)
	}
}
