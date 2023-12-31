package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const imagesDirectory = "./images"
const numImages = 30

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveRandomImage(w, r)
	})

	http.HandleFunc("/randomImages", func(w http.ResponseWriter, r *http.Request) {
		serveMultipleRandomImages(w, r, numImages)
	})

	http.ListenAndServe(":8000", nil)
}

func serveRandomImage(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(imagesDirectory)
	if err != nil {
		http.Error(w, "Unable to read images directory", http.StatusInternalServerError)
		return
	}

	var images []string
	for _, file := range files {
		if !file.IsDir() {
			images = append(images, file.Name())
		}
	}

	if len(images) == 0 {
		http.Error(w, "No images available", http.StatusNotFound)
		return
	}

	rand.Seed(time.Now().UnixNano())
	randomImage := images[rand.Intn(len(images))]

	// Set cache control headers
	w.Header().Set("Cache-Control", "no-store, must-revalidate")

	http.ServeFile(w, r, imagesDirectory+"/"+randomImage)

	http.HandleFunc("/images/", serveImage)
}

func serveMultipleRandomImages(w http.ResponseWriter, r *http.Request, numImages int) {
	var imageLinks []string

	for i := 1; i <= numImages; i++ {
		rand.Seed(time.Now().UnixNano())
		files, err := ioutil.ReadDir(imagesDirectory)
		if err != nil {
			http.Error(w, "Unable to read images directory", http.StatusInternalServerError)
			return
		}

		var images []string
		for _, file := range files {
			if !file.IsDir() {
				images = append(images, file.Name())
			}
		}

		if len(images) == 0 {
			http.Error(w, "No images available", http.StatusNotFound)
			return
		}

		randomImage := images[rand.Intn(len(images))]
		imageLinks = append(imageLinks, fmt.Sprintf("/images/%s", randomImage))
	}

	// Set cache control headers
	w.Header().Set("Cache-Control", "no-store, must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"image_links": %s}`, toJSON(imageLinks))))
}

// when call localhost:8000/images/1.png return image 1.png in folder images
func serveImage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])

}

func toJSON(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
