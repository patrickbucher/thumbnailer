package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const command = "/usr/bin/convert" // ImageMagick

func main() {
	http.HandleFunc("/canary", canary)
	http.HandleFunc("/thumbnail", thumbnail)
	log.Fatal(http.ListenAndServe("0.0.0.0:1337", nil))
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

	args, err := parseParams(r)
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

	thumbnail, err := ioutil.TempFile("", "*.png")
	defer os.Remove(thumbnail.Name())
	defer thumbnail.Close()
	if err != nil {
		log.Printf("create temp file '%s' for thumbnail: %v", thumbnail.Name(), err)
		response(w, http.StatusInternalServerError)
		return
	}

	args = append(args, "-flatten")            // white background, no transparency
	args = append(args, tempFile.Name()+"[0]") // [0] means first page
	args = append(args, thumbnail.Name())
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

// TODO: consider returning a struct with a method to create the param slice
func parseParams(r *http.Request) ([]string, error) {
	params := make([]string, 0)
	width, err := parseIntParam(r, "width")
	if err != nil {
		return params, err
	}
	height, err := parseIntParam(r, "height")
	if err != nil {
		return params, err
	}
	density, err := parseIntParam(r, "density")
	if err != nil {
		return params, err
	}
	if width > 0 || height > 0 {
		params = append(params, "-thumbnail")
		widthStr, heightStr := "", ""
		if width > 0 {
			widthStr = strconv.Itoa(width)
		}
		if height > 0 {
			heightStr = strconv.Itoa(height)
		}
		resizeParam := fmt.Sprintf("%sx%s", widthStr, heightStr)
		if widthStr != "" && heightStr != "" {
			// "!" ignores aspect ratio
			resizeParam = fmt.Sprintf("%sx%s!", widthStr, heightStr)
		}
		params = append(params, resizeParam)
	}
	if density > 0 {
		params = append(params, "-density")
		params = append(params, strconv.Itoa(density))
	}
	return params, nil
}

func parseIntParam(r *http.Request, name string) (int, error) {
	var intVal int
	stringParam := r.FormValue(name)
	if stringParam != "" {
		i, err := strconv.Atoi(stringParam)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("'%s' param not an integer: %v", name, err))
		}
		intVal = i
	}
	return intVal, nil
}
