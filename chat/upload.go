package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type uploadHandler struct{}

func (uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userID")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filename := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
