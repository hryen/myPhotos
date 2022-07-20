package web

import (
	"embed"
	"github.com/gorilla/mux"
	"io/fs"
	"myPhotos/api"
	"myPhotos/config"
	"myPhotos/logger"
	"net/http"
	"os"
)

//go:embed static
var static embed.FS

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/medias", api.ListMedia).Methods("GET")
	r.HandleFunc("/api/medias/count", api.GetMediaCount).Methods("GET")
	r.HandleFunc("/api/medias/search", api.SearchMedia).Methods("GET")
	r.HandleFunc("/api/medias/{id}", api.GetMedia).Methods("GET")
	r.HandleFunc("/api/medias/{id}/thumbnail", api.GetMediaThumbnail).Methods("GET")
	r.HandleFunc("/api/medias/{id}/video", api.GetLivePhotoVideo).Methods("GET")
	r.HandleFunc("/api/medias/{id}/info", api.GetMediaInfo).Methods("GET")

	// frontend
	staticFS, _ := fs.Sub(static, "static")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(staticFS))))

	logger.InfoLogger.Println("About to listen on " + config.HTTPPort + ", Go to http://" + config.HTTPAddr)
	err := http.ListenAndServe(config.HTTPAddr+":"+config.HTTPPort, r)
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
}
