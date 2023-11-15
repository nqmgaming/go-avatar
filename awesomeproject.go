package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveRandomImage(w, r)
	})

	http.ListenAndServe(":8000", nil)
}

func serveRandomImage(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./images")
	if err != nil {
		http.Error(w, "Unable to read images directory", http.StatusInternalServerError)
		return
	}

	// Filter out directories, only keep regular files
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

	http.ServeFile(w, r, "./images/"+randomImage)
}
