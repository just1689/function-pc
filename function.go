package F

import (
	"cloud.google.com/go/datastore"
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/plancks-cloud/function-pc/io"
	"go.opencensus.io/trace"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	config     *io.Configuration
	configFunc = defaultConfigFunc
	pwd        = os.Getenv("pwd")
	projectId  = os.Getenv("GCP_PROJECT")
)

func Handler(w http.ResponseWriter, r *http.Request) {

	action := r.URL.Query().Get("action")
	collection := r.URL.Query().Get("collection")

	if githubRepo == "" {
		io.WriteError(w, http.StatusBadRequest, "Bad request: ?project=X found.")
		return
	}

	config = &io.Configuration{}
	config.Once.Do(func() { configFunc() })

	//Get clones from Github API
	u, err := getClones(githubRepo)
	if err != nil {
		fmt.Println(err)
		io.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	//Store all the clone data
	storeObjects(githubRepo, u)

	days, err := config.ListAll(githubRepo)
	if err != nil {
		fmt.Println(err)
		io.WriteError(w, http.StatusInternalServerError, err.Error())
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
	io.WriteURLToWriter(url, w)
}

func getClones(project string) (clones *io.Clones, err error) {
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

	clones = &io.Clones{}
	err = json.Unmarshal(body, clones)
	return

}

func getDataStoreKeys(project string, clones *io.Clones) (keys []*datastore.Key) {
	keys = make([]*datastore.Key, len(clones.Days))
	for _, c := range clones.Days {
		key := datastore.NameKey(project, c.Timestamp, nil)
		keys = append(keys, key)
	}
	return keys
}

func countUniques(days []*io.Day) int {
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

func init() {
	config = &io.Configuration{}
}

func defaultConfigFunc() {

	stackdriverExporter, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: projectId})
	if err != nil {
		config.Err = err
		return
	}

	trace.RegisterExporter(stackdriverExporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	client, err := datastore.NewClient(context.Background(), projectId)
	if err != nil {
		config.Err = err
		return
	}

	config.DataStoreClient = client
}
