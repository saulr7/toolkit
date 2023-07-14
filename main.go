package main

import (
	"fmt"
	"log"
	"net/http"
	"toolkit/toolkit"
)

func main() {

	// var tools toolkit.Tools
	// s := tools.RandomString(20)
	// fmt.Println(s)

	mux := routes()

	log.Println("Starting server on port 8080")

	http.ListenAndServe(":8080", mux)
}

func routes() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/upload", uploadFiles)
	mux.HandleFunc("/upload-one", uploadOneFile)

	return mux

}

func uploadFiles(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	t := toolkit.Tools{
		MaxFileSize:      1024 * 1024 * 1024,
		AllowedFileTypes: []string{"image/jpg", "image/png", "image/gif", "image/jpeg"},
	}

	files, err := t.UploadFiles(r, "./uploads")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := "hey:\n"

	for _, item := range files {
		out += fmt.Sprintf("Uploaded %s to the uploads folder, renamed to %s\n", item.OriginalFileName, item.NewFileName)
	}

	_, _ = w.Write([]byte(out))

}

func uploadOneFile(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	t := toolkit.Tools{
		MaxFileSize:      1024 * 1024 * 1024,
		AllowedFileTypes: []string{"image/jpg", "image/png", "image/gif", "image/jpeg"},
	}

	file, err := t.UploadOneFile(r, "./uploads")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := "hey:\n"
	out += fmt.Sprintf("Uploaded %s to the uploads folder, renamed to %s\n", file.OriginalFileName, file.NewFileName)

	_, _ = w.Write([]byte(out))
}
