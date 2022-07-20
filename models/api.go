package models

import "time"

type ApiResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewApiResponse(ok bool, msg string, data interface{}) *ApiResponse {
	ar := new(ApiResponse)
	if ok {
		ar.Code = 0
		ar.Data = data
	} else {
		ar.Code = -1
	}
	ar.Msg = msg
	ar.Timestamp = time.Now().Local()
	return ar
}
