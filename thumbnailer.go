package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

const (
	command = "/usr/bin/convert" // imagmagick
	density = "-density"
	dpi     = "150"
	resize  = "-resize"
	width   = "200"
)

func main() {
	http.HandleFunc("/thumbnail", thumbnail)
	log.Fatal(http.ListenAndServe("0.0.0.0:1337", nil))
}

func thumbnail(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("pdf")
	defer file.Close()
	if err != nil {
		log.Print(err)
		return
	}

	tempFile, err := ioutil.TempFile("", "*.pdf")
	defer tempFile.Close()
	if err != nil {
		log.Print(err)
		return
	}

	io.Copy(tempFile, file)
	thumbnail, err := ioutil.TempFile("", "*.png")
	defer thumbnail.Close()
	if err != nil {
		log.Print(err)
		return
	}

	source := tempFile.Name() + "[0]" // [0] means first page
	target := thumbnail.Name()
	cmd := exec.Command(command, density, dpi, resize, width, source, target)
	err = cmd.Run()
	if err != nil {
		log.Print(err)
		return
	}

	http.ServeFile(w, r, thumbnail.Name())
}
