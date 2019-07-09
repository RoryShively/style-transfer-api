package main

import (
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

type Image struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Name string `json:"name"`
	Style string `json:"style"`
	Status string `json:"status"`
	UploadedImagePath string `json:"uploadedImagePath"`
	StylizedImagePath string `json:"stylizedImagePath"`
}

func (i *Image) GetStylizedImage() *os.File {
	filePath := i.StylizedImagePath

	img, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err) // perhaps handle this nicer
	}
	// defer img.Close()
	return img
}

func (i *Image) GetUploadedImage() *os.File {
	filePath := i.UploadedImagePath

	img, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err) // perhaps handle this nicer
	}
	// defer img.Close()
	return img
}

func (i *Image) UploadImage(file io.Reader, name string) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	imagePrefix := rand.Int()
	
	imagePath := fmt.Sprintf("/data/uploaded/%v-%v", imagePrefix, i.Name)

	// create file if not exists
	_, err = os.Stat(imagePath)
	if os.IsNotExist(err) {
		var file, err = os.Create(imagePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	err = ioutil.WriteFile(imagePath, bytes, 0777)
	if err != nil {
		log.Println("ioutil.Writefile()")
		log.Fatal(err)
	}

	i.UploadedImagePath = imagePath
}