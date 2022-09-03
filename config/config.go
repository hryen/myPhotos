package config

import (
	"github.com/spf13/viper"
	"myPhotos/logger"
	"os"
	"path/filepath"
	"time"
)

var (
	// DataSourceName 数据库连接字符串
	DataSourceName string

	// DataPath 数据存放路径，包含缩略图、数据库等文件
	DataPath string

	// ThumbnailPath 保存照片和视频缩略图的路径
	ThumbnailPath string

	// UploadPath 保存上传的照片和视频的路径
	UploadPath string

	// GPSToGeo 是否开启通过 GPS 经纬度数据使用高德地图 api 反查地理位置
	GPSToGeo bool

	// AMapKey 高德地图 key
	AMapKey string

	// PhotoExtList 照片的扩展名列表，小写
	PhotoExtList []string

	// VideoExtList 视频的扩展名列表，小写
	VideoExtList []string

	// Location 时区
	Location *time.Location

	// HTTPAddr web服务监听地址
	HTTPAddr string

	// HTTPPort web服务监听端口
	HTTPPort string
)

func InitializeConfig(path string) {
	DataPath = filepath.Join(path, "data")

	// 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		logger.ErrorLogger.Println("fatal error config file:", err)
		os.Exit(1)
	}

	DataSourceName = viper.GetString("DataSourceName")
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
		err := os.MkdirAll(ThumbnailPath, 0755)
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
		err := os.MkdirAll(UploadPath, 0755)
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
