package web

import (
	"github.com/gorilla/mux"
	"myPhotos/api"
	"myPhotos/config"
	"myPhotos/logger"
	"net/http"
	"os"
)

func StartServer() {
	logger.InfoLogger.Println("web server init...")
	r := mux.NewRouter()

	r.HandleFunc("/api/medias", api.ListMedia).Methods(http.MethodGet)
	r.HandleFunc("/api/medias", api.UploadFile).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/medias/count", api.GetMediaCount).Methods(http.MethodGet)
	r.HandleFunc("/api/medias/search", api.SearchMedia).Methods(http.MethodGet)
	r.HandleFunc("/api/medias/{id}", api.GetMedia).Methods(http.MethodGet)
	r.HandleFunc("/api/medias/{id}/thumbnail", api.GetMediaThumbnail).Methods(http.MethodGet)
	r.HandleFunc("/api/medias/{id}/video", api.GetLivePhotoVideo).Methods(http.MethodGet)
	r.HandleFunc("/api/medias/{id}/info", api.GetMediaInfo).Methods(http.MethodGet)

	r.Use(mux.CORSMethodMiddleware(r))

	logger.InfoLogger.Println("About to listen on " + config.HTTPPort + ", Go to http://" + config.HTTPAddr)
	err := http.ListenAndServe(config.HTTPAddr+":"+config.HTTPPort, r)
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
}
