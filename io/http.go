package io

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func WriteURLToWriter(url string, w http.ResponseWriter) {
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	_, err := w.Write(body)
	if err != nil {
		fmt.Println(err)
	}
}

func WriteError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Add("Content-type", "text")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("ETag", "0")
	WriteURLToWriter(msg, w)
}
