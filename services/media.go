package services

import (
	"myPhotos/config"
	"myPhotos/entity"
	"myPhotos/third_party/amap"
	"myPhotos/third_party/exiftool"
	"myPhotos/third_party/thumbnail"
	"myPhotos/tools"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func SaveMedia(metadata exiftool.FileMetadata, id string, originalFilename string) {
	m := convertToMedia(metadata)
	m.OriginalFilename = originalFilename

	// 获取地理位置
	if config.GPSToGeo && m.GPSLongitude != "" && m.GPSLatitude != "" {
		g := amap.GetMediaGeo(m.GPSLongitude, m.GPSLatitude)
		m.GPSCountry = g.Country
		m.GPSCity = g.City
		m.GPSProvince = g.Province
		m.GPSDistrict = g.District
		m.GPSFormattedAddress = g.FormattedAddress
	}

	m.ID = id

	// 生成缩略图，只有动态图片的视频不生成
	if m.MediaType != entity.LivePhotoVideo {
		thumbnail.SaveMediaThumbnail(&m)
	}

	// 插入数据库
	entity.SaveMedia(m)
}

func convertToMedia(metadata exiftool.FileMetadata) entity.Media {
	path := metadata.GetString("SourceFile")
	m := entity.Media{Path: path}

	ext := strings.ToLower(filepath.Ext(path))

	// photo
	if tools.ArrayContains(config.PhotoExtList, ext) {
		m.MediaType = entity.Photo

		if metadata.GetString("DateTimeOriginal") == "" {
			m.DateTime = tools.GetModTime(path)
		} else {
			times := strings.Split(metadata.GetString("DateTimeOriginal"), " ")
			time1 := strings.Split(times[1], ":")
			if i, _ := strconv.Atoi(time1[0]); i > 23 {
				time1[0] = "00"
			}
			if i, _ := strconv.Atoi(time1[1]); i > 59 {
				time1[1] = "00"
			}
			if i, _ := strconv.Atoi(time1[2]); i > 59 {
				time1[2] = "00"
			}
			d := times[0] + " " + strings.Join(time1, ":")
			t, _ := time.ParseInLocation("2006:01:02 15:04:05", d, config.Location)
			m.DateTime = t
		}

		m.MediaGroupUUID = metadata.GetString("MediaGroupUUID")
		if m.MediaGroupUUID != "" {
			m.MediaType = entity.LivePhoto
		}
	}

	// video
	if tools.ArrayContains(config.VideoExtList, ext) {
		m.MediaType = entity.Video

		t := ""
		createDate := metadata.GetString("CreateDate")
		creationDate := metadata.GetString("CreationDate")
		if createDate != "" && createDate != "0000:00:00 00:00:00" {
			t = createDate
		} else if creationDate != "" && creationDate != "0000:00:00 00:00:00" {
			t = creationDate
		}
		if t != "" {
			t, _ := time.ParseInLocation("2006:01:02 15:04:05", t, config.Location)
			t = t.Add(time.Hour * 8)
			m.DateTime = t
		} else {
			m.DateTime = tools.GetModTime(path)
		}

		m.Duration = tools.ConvertDuration(metadata.GetFloat("Duration"))

		m.ContentIdentifier = metadata.GetString("ContentIdentifier")
		if m.ContentIdentifier != "" {
			m.MediaType = entity.LivePhotoVideo
		}
	}

	m.FileType = metadata.GetString("FileType")
	m.FileSize = metadata.GetString("FileSize")
	m.ImageSize = metadata.GetString("ImageSize")
	m.Make = metadata.GetString("Make")
	m.Model = metadata.GetString("Model")
	m.Megapixels = metadata.GetString("Megapixels")
	m.GPSLongitude = metadata.GetString("GPSLongitude")
	m.GPSLatitude = metadata.GetString("GPSLatitude")

	m.ISO = metadata.GetString("ISO")
	m.Flash = metadata.GetString("Flash")
	m.FocalLength = metadata.GetString("FocalLength")
	m.ShutterSpeed = metadata.GetString("ShutterSpeed")
	m.Aperture = metadata.GetString("Aperture")

	return m
}
