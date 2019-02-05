package F

import (
	"github.com/plancks-cloud/function-pc/io"
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
		handleSet()
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
	} else if collection == io.ServiceCollectionName {
		sl, err := io.ListAllServices(config)
		if err != nil {
			io.WriteError(w, http.StatusInternalServerError, "Could not list all services")
			return
		}
		io.WriteObjectToJson(w, sl)
	} else {
		io.WriteError(w, http.StatusBadRequest, "Bad request: action must be get or set.")
		return
	}

}

func handleSet(collection string, config io.Configuration, id string) {
	if collection == io.RouteCollectionName {

		io.StoreRoutes(config, id)
	}
}

func init() {
	config = &io.Configuration{}
}
