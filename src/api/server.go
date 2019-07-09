package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CreateServer() *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/image/{id:[0-9]+}",  getImageData).
		Methods("GET")
	r.HandleFunc("/image",  postImage).
		Methods("POST")
	r.HandleFunc("/image/stylized/{id:[0-9]+}",  getStylizedImage).
		Methods("GET")
	r.HandleFunc("/image/uploaded/{id:[0-9]+}",  getUploadedImage).
		Methods("GET")
	r.HandleFunc("/models",  getAvailableModels).
		Methods("GET")
	

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API"))
	})

	return &http.Server{
        Addr:         "0.0.0.0:3100",
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: r, // Pass our instance of gorilla/mux in.
	}
}