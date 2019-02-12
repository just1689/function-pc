package io

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

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
