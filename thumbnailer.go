package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/patrickbucher/thumbnailer"
)

const command = "/usr/bin/convert" // ImageMagick

func main() {
	http.HandleFunc("/canary", canary)
	http.HandleFunc("/thumbnail", thumbnail)
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe("0.0.0.0:"+PORT, nil))
}

func canary(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("beep, beep\n"))
}

func thumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("called with method %s", r.Method)
		response(w, http.StatusMethodNotAllowed)
		return
	}

	params, err := thumbnailer.ParseParams(r)
	if err != nil {
		log.Printf("parsing params from multipart request: %v", err)
		response(w, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("pdf")
	defer file.Close()
	if err != nil {
		log.Printf("get 'pdf' file from multipart request: %v", err)
		response(w, http.StatusBadRequest)
		return
	}

	tempFile, err := ioutil.TempFile("", "*.pdf")
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	if err != nil {
		log.Printf("create temp file '%s' for source PDF: %v", err, tempFile.Name())
		response(w, http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Printf("copy request file to '%s' failed: %v", tempFile.Name(), err)
		response(w, http.StatusInternalServerError)
		return
	}

	thumbnail, err := ioutil.TempFile("", "*."+params.Format)
	defer os.Remove(thumbnail.Name())
	defer thumbnail.Close()
	if err != nil {
		log.Printf("create temp file '%s' for thumbnail: %v", thumbnail.Name(), err)
		response(w, http.StatusInternalServerError)
		return
	}

	inputArg := tempFile.Name() + "[0]" // [0] means first page
	args := params.AsArgs("-flatten", "-strip", inputArg, thumbnail.Name())
	cmd := exec.Command(command, args...)

	started := time.Now()
	err = cmd.Run()
	finished := time.Now()
	if err != nil {
		log.Printf("executing '%s %v': %v", command, strings.Join(args, " "), err)
		response(w, http.StatusInternalServerError)
		return
	}
	duration := finished.Sub(started)
	http.ServeFile(w, r, thumbnail.Name())
	log.Printf("%s %s [%v]", command, strings.Join(args, " "), duration)
}

func response(w http.ResponseWriter, statusCode int) {
	statusMessage := http.StatusText(statusCode)
	w.WriteHeader(statusCode)
	w.Write([]byte(statusMessage))
}
