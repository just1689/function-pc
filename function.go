package F

import (
	"encoding/json"
	"github.com/plancks-cloud/function-pc/io"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	config     *io.Configuration
	configFunc = io.DefaultConfigFunc
	projectId  = os.Getenv("GCP_PROJECT")
)

func Handler(w http.ResponseWriter, r *http.Request) {

	action := r.URL.Query().Get("action")
	collection := r.URL.Query().Get("collection")

	if action == "" || collection == "" {
		io.WriteError(w, http.StatusBadRequest, "Bad request: action and collection required.")
		return
	}

	config = &io.Configuration{}
	config.Once.Do(func() { configFunc(projectId, config) })

	id := r.Header.Get("persist-id")
	key := r.Header.Get("persist-key")

	authenticated := io.Auth(id, key)
	if !authenticated {
		io.WriteError(w, http.StatusUnauthorized, "Unauthorized to access. Check ID and key.")
		return
	}

	if action == "get" {
		handleGet(collection, w)
	} else if action == "set" {
		handleSet(collection, config, id, w, r)
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

func handleSet(collection string, config *io.Configuration, id string, w http.ResponseWriter, r *http.Request) {
	if collection == io.RouteCollectionName {
		var routes []io.Route
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&routes)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not decode routes")
			return
		}
		err = io.StoreRoutes(config, id, routes)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not store routes")
			return
		}
		io.WriteObjectToJson(w, "")
	} else if collection == io.ServiceCollectionName {
		var sl []io.Service
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&sl)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not decode services")
			return
		}
		err = io.StoreServices(config, id, sl)
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
