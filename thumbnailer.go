package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	command = "/usr/bin/convert" // ImageMagick
	density = "-density"
	dpi     = "150"
	resize  = "-resize"
	width   = "400"
)

func main() {
	http.HandleFunc("/thumbnail", thumbnail)
	log.Fatal(http.ListenAndServe("0.0.0.0:1337", nil))
}

// TODO: accept density and resolution as parameters
func thumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("called with method %s", r.Method)
		response(w, http.StatusMethodNotAllowed)
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

	thumbnail, err := ioutil.TempFile("", "*.png")
	defer os.Remove(thumbnail.Name())
	defer thumbnail.Close()
	if err != nil {
		log.Printf("create temp file '%s' for thumbnail: %v", thumbnail.Name(), err)
		response(w, http.StatusInternalServerError)
		return
	}

	source := tempFile.Name() + "[0]" // [0] means first page
	target := thumbnail.Name()
	args := []string{density, dpi, resize, width, source, target}
	cmd := exec.Command(command, args...)
	err = cmd.Run()
	if err != nil {
		log.Printf("executing '%s %v': %v", command, strings.Join(args, " "), err)
		response(w, http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, thumbnail.Name())
}

func response(w http.ResponseWriter, statusCode int) {
	statusMessage := http.StatusText(statusCode)
	w.WriteHeader(statusCode)
	w.Write([]byte(statusMessage))
}
