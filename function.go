package F

import (
	"encoding/json"
	"github.com/plancks-cloud/function-pc/controller"
	"github.com/plancks-cloud/function-pc/domain"
	"github.com/plancks-cloud/function-pc/io"
	"github.com/sirupsen/logrus"
	iog "io"
	"net/http"
	"os"
)

var (
	config     = &io.Configuration{}
	configFunc = io.DefaultConfigFunc
	projectId  = os.Getenv("GCP_PROJECT")
)

type requestDescription struct {
	method     string
	collection string
	id         string
	key        string
	body       iog.ReadCloser
}

func Handler(w http.ResponseWriter, r *http.Request) {

	req := describeRequest(r)

	if req.method == "" || req.collection == "" {
		io.WriteError(w, http.StatusBadRequest, "Bad request: action and collection required.")
		return
	}

	config.Once.Do(func() { configFunc(projectId, config) })

	authenticated := io.Auth(req.id, req.key)
	if !authenticated {
		io.WriteError(w, http.StatusUnauthorized, "Unauthorized to access. Check ID and key.")
		return
	}

	if req.method == http.MethodGet {
		handleGet(req.collection, w)
	} else if req.method == http.MethodPost {
		handleSet(req, w)
	} else {
		io.WriteError(w, http.StatusBadRequest, "Bad request: action must be get or set.")
		return
	}

}

func handleGet(collection string, w http.ResponseWriter) {
	if collection == domain.RouteCollectionName {
		sl, err := controller.ListAllRoutes(config)
		if err != nil {
			io.WriteError(w, http.StatusInternalServerError, "Could not list all routes")
			return
		}
		io.WriteObjectToJson(w, sl)
		return
	} else if collection == domain.ServiceCollectionName {
		sl, err := controller.ListAllServices(config)
		if err != nil {
			io.WriteError(w, http.StatusInternalServerError, "Could not list all services")
			return
		}
		io.WriteObjectToJson(w, sl)
		return
	}
	io.WriteError(w, http.StatusBadRequest, "Bad request: action must be get or set.")

}

func handleSet(req *requestDescription, w http.ResponseWriter) {
	if req.collection == domain.RouteCollectionName {
		var routes []domain.Route
		decoder := json.NewDecoder(req.body)
		err := decoder.Decode(&routes)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not decode routes")
			return
		}
		err = controller.StoreRoutes(config, req.id, routes)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not store routes")
			return
		}
		io.WriteObjectToJson(w, "")
	} else if req.collection == domain.ServiceCollectionName {
		var sl []domain.Service
		decoder := json.NewDecoder(req.body)
		err := decoder.Decode(&sl)
		if err != nil {
			logrus.Error(err)
			io.WriteError(w, http.StatusInternalServerError, "Could not decode services")
			return
		}
		err = controller.StoreServices(config, req.id, sl)
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
		method:     r.Method,
		collection: r.URL.Query().Get("collection"),
		body:       r.Body,
	}

}
