package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func (api *API) ImgProfileView(w http.ResponseWriter, r *http.Request) {
	// View with response image `img-avatar.png` from path `assets/images`
	// TODO: answer here

	imageName := r.URL.Query().Get("avatar")                                        // mengambil nama image dari query url
	fileBytes, err := ioutil.ReadFile("./assets/images/img-avatar.png" + imageName) // membaca file image menjadi bytes
	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		w.Write([]byte("File not found"))
		return
	}
	w.WriteHeader(http.StatusOK) // 200
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes) // menampilkan image sebagai response

}

func (api *API) ImgDashboard(w http.ResponseWriter, r *http.Request) {
	// View with response image `img-avatar.png` from path `assets/images`
	// TODO: answer here

	// mengambil nama image dari query url
	fileBytes, err := ioutil.ReadFile("./assets/images/background_image_Food_delivery_man_riding_motorcycles_cartoon.png") // membaca file image menjadi bytes
	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		w.Write([]byte("File not found"))
		return
	}
	w.WriteHeader(http.StatusOK) // 200
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes) // menampilkan image sebagai response

}

func (api *API) ImgProfileUpdate(w http.ResponseWriter, r *http.Request) {
	// Update image `img-avatar.png` from path `assets/images`
	// TODO: answer here

	// check response body
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alias := "img-avatar"

	uploadedFile, handler, err := r.FormFile("file-avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	//mengambil relative path dari proyek
	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// //membuat nama file baru yang akan disimpan
	filename := handler.Filename
	if alias != "" { //kalau alias disi maka nama file = alias
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	}
	// filename := fmt.Sprintf("%s", filepath.Ext(handler.Filename))

	//membentuk lokasi tempat menyimpan file
	fileLocation := filepath.Join(dir, "assets/images", filename)
	fmt.Println("fileLocation", fileLocation)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	//mengisi file baru dengan data dari file yang ter-upload
	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

	api.dashboardView(w, r)
}
