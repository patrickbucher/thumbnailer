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
)

const command = "/usr/bin/convert" // ImageMagick

func main() {
	http.HandleFunc("/thumbnail", thumbnail)
	log.Fatal(http.ListenAndServe("0.0.0.0:1337", nil))
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

	source := tempFile.Name() + "[0]" // [0] means first page
	target := thumbnail.Name()
	args = append(args, source)
	args = append(args, target)
	cmd := exec.Command(command, args...)
	err = cmd.Run()
	if err != nil {
		log.Printf("executing '%s %v': %v", command, strings.Join(args, " "), err)
		response(w, http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, thumbnail.Name())
	log.Printf("%s %s", command, strings.Join(args, " "))
}

func response(w http.ResponseWriter, statusCode int) {
	statusMessage := http.StatusText(statusCode)
	w.WriteHeader(statusCode)
	w.Write([]byte(statusMessage))
}

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
		params = append(params, "-resize")
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
