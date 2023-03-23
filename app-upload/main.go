package main

import (
	"fmt"
	"github.com/faizanaryan94/toolkit"
	"log"
	"net/http"
)

func main() {
	mux := routes()
	log.Println("starting server on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/upload", uploadFiles)
	mux.HandleFunc("/upload-one", uploadFile)
	return mux
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tools := toolkit.Tools{
		MaxFileSize:      1024 * 1024 * 1024,
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif"},
	}
	files, err := tools.UploadFiles(r, "./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	out := ""
	for _, item := range files {
		out += fmt.Sprintf("Uploaded %s to the uploads folder, renamed to %s", item.OriginalFileName, item.NewFileName)
	}

	_, _ = w.Write([]byte(out))
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tools := toolkit.Tools{
		MaxFileSize:      1024 * 1024 * 1024,
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif"},
	}
	file, err := tools.UploadFile(r, "./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf("uploaded 1 file, %s, to the uploads folder", file.OriginalFileName)))

}
