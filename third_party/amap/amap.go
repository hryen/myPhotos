package amap

import (
	"encoding/json"
	"io"
	"myPhotos/config"
	"myPhotos/logger"
	"net/http"
)

type amapRegeoResponse struct {
	Status    string    `json:"status"`
	Regeocode regeocode `json:"regeocode"`
	Info      string    `json:"info"`
}
type regeocode struct {
	AddressComponent addressComponent `json:"addressComponent"`
	FormattedAddress string           `json:"formatted_Address"`
}
type addressComponent struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	// City 正常情况是 string, 但直辖市会返回`[]`, 直接判断类型是否是 string 即可
	City interface{} `json:"city"`
	// District 正常情况是 string, 但直辖市会返回`[]`, 直接判断类型是否是 string 即可
	District interface{} `json:"district"`
}
type GeoInfo struct {
	Country          string
	Province         string
	City             string
	District         string
	FormattedAddress string
}

const errorMsgPrefix = "failed to get media geography info:"

// GetMediaGeo 利用高德 api 通过经纬度获取地理位置 TODO 简化代码 json unmarshal 时用 map
func GetMediaGeo(lon, lat string) GeoInfo {
	resp, err := http.Get("https://restapi.amap.com/v3/geocode/regeo?key=" + config.AMapKey + "&location=" + lon + "," + lat)
	if err != nil {
		logger.ErrorLogger.Println(errorMsgPrefix, err)
		return GeoInfo{}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Println(errorMsgPrefix, err)
		return GeoInfo{}
	}

	var a amapRegeoResponse
	err = json.Unmarshal(body, &a)
	if err != nil {
		logger.ErrorLogger.Println(errorMsgPrefix, "lon:"+lon, "lat:"+lat, err)
		return GeoInfo{}
	}

	if a.Status == "0" {
		logger.ErrorLogger.Println(errorMsgPrefix, a.Info)
		return GeoInfo{}
	} else {
		city := ""
		if _, e := a.Regeocode.AddressComponent.City.(string); e {
			city = a.Regeocode.AddressComponent.City.(string)
		}
		district := ""
		if _, e := a.Regeocode.AddressComponent.District.(string); e {
			district = a.Regeocode.AddressComponent.District.(string)
		}

		return GeoInfo{
			Country:          a.Regeocode.AddressComponent.Country,
			City:             city,
			Province:         a.Regeocode.AddressComponent.Province,
			District:         district,
			FormattedAddress: a.Regeocode.FormattedAddress,
		}
	}
}
