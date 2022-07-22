package entity

import (
	"github.com/google/uuid"
	"myPhotos/logger"
	"time"
)

const Photo = "photo"
const Video = "video"
const LivePhoto = "LivePhoto"
const LivePhotoVideo = "LivePhotoVideo"

type Media struct {
	ID string `gorm:"primarykey"`

	// basic info
	MediaType  string
	Path       string
	FileType   string
	FileSize   string
	DateTime   time.Time
	ImageSize  string
	Make       string
	Model      string
	Megapixels string

	GPSLongitude        string
	GPSLatitude         string
	GPSCountry          string
	GPSProvince         string
	GPSCity             string
	GPSDistrict         string
	GPSFormattedAddress string

	// photo info
	MediaGroupUUID string
	ISO            string
	Flash          string
	FocalLength    string
	ShutterSpeed   string
	Aperture       string

	// video info
	Duration          string
	ContentIdentifier string
}

func SaveMedia(m Media) {
	m.ID = uuid.New().String()
	err := DB.Create(&m).Error
	if err != nil {
		logger.ErrorLogger.Println("save media error:", err.Error())
	} else {
		logger.InfoLogger.Println("save media:", m.Path)
	}
}

func DeleteMediaByPath(path string) {
	var m Media
	err := DB.First(&m, "path = ?", path).Error
	if err != nil {
		logger.ErrorLogger.Println("delete media error:", err)
	} else {
		err = DB.Delete(&m).Error
		if err != nil {
			logger.ErrorLogger.Println("delete media error:", err)
		} else {
			logger.InfoLogger.Println("delete media: " + path)
		}
	}
}
