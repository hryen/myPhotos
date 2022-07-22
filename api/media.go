package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"myPhotos/config"
	"myPhotos/entity"
	"myPhotos/logger"
	"myPhotos/models"
	"net/http"
	"path/filepath"
	"strconv"
)

// ListMedia 分页返回 Media，按时间排序，排除动态图片的视频
func ListMedia(w http.ResponseWriter, r *http.Request) {
	var m []models.ApiMedia
	err := entity.DB.Model(&entity.Media{}).Scopes(withoutLivePhotoVideo).Order("date_time desc").Scopes(paginate(r)).Find(&m).Error
	if err != nil {
		writeJSON(w, models.NewApiResponse(false, "failed to list medias: "+err.Error(), nil))
	} else {
		writeJSON(w, models.NewApiResponse(true, "list the medias successfully", m))
	}
}

// GetMediaCount 获取所有媒体的数量，排除动态图片的视频
func GetMediaCount(w http.ResponseWriter, _ *http.Request) {
	var totalCount int64
	err := entity.DB.Model(&entity.Media{}).Scopes(withoutLivePhotoVideo).Count(&totalCount).Error
	if err != nil {
		writeJSON(w, models.NewApiResponse(false, "failed to count the medias:"+err.Error(), nil))
	} else {
		writeJSON(w, models.NewApiResponse(true, "count the medias successfully", totalCount))
	}
}

func GetMedia(w http.ResponseWriter, r *http.Request) {
	doGetMedia(w, r, false)
}

func GetMediaThumbnail(w http.ResponseWriter, r *http.Request) {
	doGetMedia(w, r, true)
}

func doGetMedia(w http.ResponseWriter, r *http.Request, isThumbnail bool) {
	id := mux.Vars(r)["id"]
	var m entity.Media

	err := entity.DB.First(&m, "id = ?", id).Error
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		logger.ErrorLogger.Println(err)
		return
	}

	if isThumbnail {
		m.Path = filepath.Join(config.ThumbnailPath, m.ID+".jpg")
	}
	// 缓存30天
	//w.Header().Set("Cache-Control", "max-age=2592000")
	http.ServeFile(w, r, m.Path)
}

func GetLivePhotoVideo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var m entity.Media

	// Get LivePhoto
	err := entity.DB.First(&m, "id = ?", id).Error
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		logger.ErrorLogger.Println(err)
		return
	}

	// Get LivePhotoVideo
	var v entity.Media
	err = entity.DB.Where("media_type = ? AND content_identifier = ?", entity.LivePhotoVideo, m.MediaGroupUUID).Find(&v).Error
	if err != nil || v.Path == "" {
		http.Error(w, "500", http.StatusInternalServerError)
		logger.ErrorLogger.Println("live photo video query error:", err)
		return
	}

	// 30天
	//w.Header().Set("Cache-Control", "max-age=2592000")
	http.ServeFile(w, r, v.Path)
}

func SearchMedia(w http.ResponseWriter, r *http.Request) {
	var ms []models.ApiMedia
	values := r.URL.Query()
	mediaType := values.Get("media_type")
	deviceMake := values.Get("make")
	deviceModel := values.Get("model")
	country := values.Get("country")
	province := values.Get("province")
	city := values.Get("city")
	district := values.Get("district")
	address := values.Get("address")
	after := values.Get("after")
	before := values.Get("before")

	rawSql := "SELECT id, media_type, date_time, duration FROM media WHERE media_type != 'LivePhotoVideo'"
	if mediaType != "" {
		rawSql += " AND media_type = '" + mediaType + "'"
	}
	if deviceMake != "" {
		rawSql += " AND make LIKE '%" + deviceMake + "%'"
	}
	if deviceModel != "" {
		rawSql += " AND model LIKE '%" + deviceModel + "%'"
	}

	if country != "" {
		rawSql += " AND gps_country LIKE '%" + country + "%'"
	}
	if province != "" {
		rawSql += " AND gps_province LIKE '%" + province + "%'"
	}
	if city != "" {
		rawSql += " AND gps_city LIKE '%" + city + "%'"
	}
	if district != "" {
		rawSql += " AND gps_district LIKE '%" + district + "%'"
	}
	if address != "" {
		rawSql += " AND (gps_formatted_address LIKE '%" + address + "%' OR gps_country LIKE '%" + address + "%')"
	}

	if after != "" {
		rawSql += " AND date_time > '" + after + "'"
	}
	if before != "" {
		before += " 00:00:00"
		rawSql += " AND date_time < '" + before + "'"
	}
	rawSql += " ORDER BY date_time DESC"
	//fmt.Println(rawSql)

	err := entity.DB.Raw(rawSql).Find(&ms).Error
	if err != nil {
		writeJSON(w, models.NewApiResponse(false, "failed to search medias: "+err.Error(), nil))
		logger.ErrorLogger.Println(err)
		return
	} else {
		writeJSON(w, models.NewApiResponse(true, "search the medias successfully", ms))
	}
}

func GetMediaInfo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var m entity.Media

	err := entity.DB.First(&m, "id = ?", id).Error
	if err != nil {
		writeJSON(w, models.NewApiResponse(false, "failed to get media: "+err.Error(), nil))
	} else {
		writeJSON(w, models.NewApiResponse(true, "get the media successfully", m))
	}
}

func withoutLivePhotoVideo(db *gorm.DB) *gorm.DB {
	return db.Not("media_type = ?", entity.LivePhotoVideo)
}

func writeJSON(w http.ResponseWriter, v interface{}) {

	j, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// TODO cors
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(j)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
}

func paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()

		page, err := strconv.Atoi(q.Get("page"))
		if err != nil || page < 1 {
			_ = db.AddError(fmt.Errorf("page cannot be less than one"))
			return db
		}

		pageSize, err := strconv.Atoi(q.Get("page_size"))
		if err != nil || pageSize < 1 {
			_ = db.AddError(fmt.Errorf("page_size cannot be less than one"))
			return db
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
