package main

import (
	"encoding/json"
	"log"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func getImageData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	image := Image{}
	dbmgr.DB.First(&image, id)

	w.Header().Set("Content-Type", "application/json")

	if image.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(image)
	}
}

func getAvailableModels(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("/models")
    if err != nil {
        log.Fatal(err)
    }

	var fileNames []string
    for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	
	var names []string
	for _, f := range fileNames {
		if strings.Contains(f, ".tar.gz") {
			continue
		} else {
			names = append(names, f[6:])
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

// Fetch image endpoint
// This will be GET of /image/:id
func getStylizedImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	image := &Image{}
	dbmgr.DB.First(image, id)
	
	file := image.GetStylizedImage()
	// defer img.Close()

	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(w, file)
}

// Fetch image endpoint
// This will be GET of /image/:id
func getUploadedImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	image := &Image{}
	dbmgr.DB.First(image, id)
	
	file := image.GetStylizedImage()
	// defer img.Close()

	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(w, file)
}

// Upload endpoint
func postImage(w http.ResponseWriter, r *http.Request) {
	// Parse form to get file
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Println("Handle this with 500 or something")
		log.Fatal(err)
	}
	defer file.Close()

	image := &Image{
		Name: handler.Filename,
		Style: r.FormValue("style"),
		Status: "uploaded",
	}
	image.UploadImage(file, handler.Filename)

	dbmgr.DB.Create(&image)

	// Publish event in redis
	if err := pubsub.Client.Publish("imageEvents", image.ID).Err(); err != nil {
		log.Fatal(err)
	}


	// Return Image JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}