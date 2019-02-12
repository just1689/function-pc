package io

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func WriteError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Add("Content-type", "text")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("ETag", "0")
	_, err := w.Write([]byte(msg))
	if err != nil {
		logrus.Error(err)
	}

}

func WriteObjectToJson(w http.ResponseWriter, object interface{}) (err error) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")
	w.Header().Add("Cache-Control", "no-cache")
	b, err := json.Marshal(object)
	if err != nil {
		logrus.Error(err)
		return
	}
	_, err = w.Write(b)
	return
}
