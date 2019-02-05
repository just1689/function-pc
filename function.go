package F

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/datastore"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

var (
	config     *configuration
	configFunc = defaultConfigFunc
	pwd        = os.Getenv("pwd")
	projectId  = os.Getenv("GCP_PROJECT")
)

type configuration struct {
	datastoreClient *datastore.Client
	err             error
	once            sync.Once
}

type Clones struct {
	Uniques int   `json:"uniques"`
	Days    []Day `json:"clones"`
}

type Day struct {
	Timestamp string `json:"timestamp"`
	Uniques   int    `json:"uniques"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	githubRepo := r.URL.Query().Get("project")
	if githubRepo == "" {
		writeError(w, http.StatusBadRequest, "Bad request: ?project=X found.")
		return
	}

	config = &configuration{}
	config.once.Do(func() { configFunc() })

	//Get clones from Github API
	u, err := getClones(githubRepo)
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	//Store all the clone data
	storeObjects(githubRepo, u)

	days, err := config.listAll(githubRepo)
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Total unique clones
	total := countUniques(days)

	//Create an url for the shield
	url := fmt.Sprintf("https://img.shields.io/badge/Unique%sclones-%v-brightgreen.svg", " ", total)

	//Hash the url
	sha := shaStr(url)

	w.Header().Add("Content-type", "image/svg+xml;charset=utf-8")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("ETag", sha)
	writeURLToWriter(url, w)
}

func getClones(project string) (clones *Clones, err error) {
	fmt.Println("Looking up clones on Github for project: ", project)
	url := fmt.Sprintf("https://api.github.com/repos/%v/traffic/clones", project)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", pwd)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "102f0b23-c735-4d04-bb80-c1e099502d37")

	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	clones = &Clones{}
	err = json.Unmarshal(body, clones)
	return

}

func storeObjects(project string, clones *Clones) {
	keys := getDataStoreKeys(project, clones)
	_, err := config.datastoreClient.PutMulti(context.Background(), keys, &clones.Days)
	if err != nil {
		fmt.Println(err)
	}
}

func getDataStoreKeys(project string, clones *Clones) (keys []*datastore.Key) {
	keys = make([]*datastore.Key, len(clones.Days))
	for _, c := range clones.Days {
		key := datastore.NameKey(project, c.Timestamp, nil)
		keys = append(keys, key)
	}
	return keys
}

func writeURLToWriter(url string, w http.ResponseWriter) {
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	_, err := w.Write(body)
	if err != nil {
		fmt.Println(err)
	}

}

func (db *configuration) listAll(project string) ([]*Day, error) {
	ctx := context.Background()
	sl := make([]*Day, 0)
	q := datastore.NewQuery(project)
	_, err := db.datastoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list: %v", err)
	}
	return sl, nil
}

func countUniques(days []*Day) int {
	total := 0
	for _, d := range days {
		total += d.Uniques
	}
	return total
}

func shaStr(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Add("Content-type", "text")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("ETag", "0")
	writeURLToWriter(msg, w)

}

func init() {
	config = &configuration{}
}

func defaultConfigFunc() {

	stackdriverExporter, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: projectId})
	if err != nil {
		config.err = err
		return
	}

	trace.RegisterExporter(stackdriverExporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	client, err := datastore.NewClient(context.Background(), projectId)
	if err != nil {
		config.err = err
		return
	}

	config.datastoreClient = client
}
