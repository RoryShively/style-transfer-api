// +build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	endpoint string
	imgPath string
	imgName string
	imgStyle string
)

func setup() {
	endpoint = os.Getenv("API_URL")
	imgPath = os.Getenv("TEST_IMAGE_PATH")
	imgStyle = "van-gogh"

	_imgPathSlice := strings.Split(imgPath, "/")
	imgName = _imgPathSlice[len(_imgPathSlice)-1]

	if endpoint == "" {
		fmt.Println("API_URL must be set")
		os.Exit(1)
	}

	if imgPath == "" {
		fmt.Println("TEST_IMAGE_PATH must be set")
		os.Exit(1)
	}
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func PostImage() (*http.Response, error) {
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	fileWriter, err := w.CreateFormFile("image", imgName)
	if err != nil {
		return nil, err
	}

	fh, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}

	w.WriteField("style", imgStyle)

	contentType := w.FormDataContentType()
	w.Close()

	resp, err := http.Post(fmt.Sprintf("%v/image", endpoint), contentType, b)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func TestUpload(t *testing.T) {
	setup()

	image := &Image{}

	t.Run("UploadImage", func(t *testing.T) {

		resp, err := PostImage()
		if err != nil {
			fmt.Println(resp.Status)
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Logf(resp.Status)
			t.Errorf("Image could not be uploaded")
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(body, image)

		if image.Status != "uploaded" {
			t.Errorf("Image does not have status: uploaded")
		}
		if image.Name != imgName {
			t.Errorf("Image does not have name: %v", imgName)
		}
		if image.Style != imgStyle {
			t.Errorf("Image does not have style: %v", imgStyle)
		}
	})

	fmt.Println("Polling image till style transfer is complete. Timeout is set at 120 seconds")

	t.Run("ImagePolling", func(t *testing.T) {

		limit := time.After(120 * time.Second)
		ticker := time.Tick(5 * time.Second)
		timeout := false
		
		   
		for image.Status == "uploaded" && !timeout {
			select {
			// Got a timeout! fail with a timeout error
			case <-limit:
				timeout = true
				t.Errorf("Test timed out")
			// Got a tick, we should check on checkSomething()
			case <-ticker:
				resp, err := http.Get(fmt.Sprintf("%v/image/%v", endpoint, image.ID))
				if err != nil {
					t.Errorf("API errored out: %v", err.Error())
				}
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Image polling Failed")
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				json.Unmarshal(body, image)
			}
		}

		if image.Status == "failed" {
			t.Errorf("Image failed during stylization")
		}

		if image.Status == "done" {
			t.Log("Image successfully stylized")
		}

	})
}
