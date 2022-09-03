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

	// frontend
	staticFS, _ := fs.Sub(static, "static")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(staticFS))))

	r.Use(mux.CORSMethodMiddleware(r))

	addr := config.HTTPAddr + ":" + config.HTTPPort
	logger.InfoLogger.Printf("Server started at: http://%s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
}
