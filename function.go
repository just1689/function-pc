package F

import (
	"encoding/json"
	"github.com/plancks-cloud/function-pc/io"
	"github.com/sirupsen/logrus"
	iog "io"
	"net/http"
	"os"
)

var (
	config     *io.Configuration
	configFunc = io.DefaultConfigFunc
	projectId  = os.Getenv("GCP_PROJECT")
)

type requestDescription struct {
	action     string
	collection string
	id         string
	key        string
	body       iog.ReadCloser
}

func Handler(w http.ResponseWriter, r *http.Request) {

	req := describeRequest(r)

	if req.action == "" || req.collection == "" {
		io.WriteError(w, http.StatusBadRequest, "Bad request: action and collection required.")
		return
	}

	config = &io.Configuration{}
	config.Once.Do(func() { configFunc(projectId, config) })

	authenticated := io.Auth(req.id, req.key)
	if !authenticated {
		io.WriteError(w, http.StatusUnauthorized, "Unauthorized to access. Check ID and key.")
		return
	}

	if req.action == "get" {
		handleGet(req.collection, w)
	} else if req.action == "set" {
		handleSet(req, config, w)
	} else {
		io.WriteError(w, http.StatusBadRequest, "Bad request: action must be get or set.")
		return
	}

}

func handleGet(collection string, w http.ResponseWriter) {
	if collection == io.RouteCollectionName {
		sl, err := io.ListAllRoutes(config)
		if err != nil {
			io.WriteError(w, http.StatusInternalServerError, "Could not list all routes")
			return
		}
		io.WriteObjectToJson(w, sl)
		return
	} else if collection == io.ServiceCollectionName {
		sl, err := io.ListAllServices(config)
		if err != nil {
			io.WriteError(w, http.StatusInternalServerError, "Could not list all services")
			return
		}
		io.WriteObjectToJson(w, sl)
		return
	}
	io.WriteError(w, http.StatusBadRequest, "Bad request: action must be get or set.")

}

func handleSet(req *requestDescription, config *io.Configuration, w http.ResponseWriter) {
	if req.collection == io.RouteCollectionName {
		var routes []io.Route
		decoder := json.NewDecoder(req.body)
		err := decoder.Decode(&routes)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not decode routes")
			return
		}
		err = io.StoreRoutes(config, req.id, routes)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not store routes")
			return
		}
		io.WriteObjectToJson(w, "")
	} else if req.collection == io.ServiceCollectionName {
		var sl []io.Service
		decoder := json.NewDecoder(req.body)
		err := decoder.Decode(&sl)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not decode services")
			return
		}
		err = io.StoreServices(config, req.id, sl)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not store services")
			return
		}
		io.WriteObjectToJson(w, "")

	}
}

func init() {
	config = &io.Configuration{}
}

func describeRequest(r *http.Request) *requestDescription {
	return &requestDescription{
		id:         r.Header.Get("persist-id"),
		key:        r.Header.Get("persist-key"),
		action:     r.URL.Query().Get("action"),
		collection: r.URL.Query().Get("collection"),
		body:       r.Body,
	}

}
