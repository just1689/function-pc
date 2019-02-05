package io

import (
	"cloud.google.com/go/datastore"
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
	"sync"
)

type Configuration struct {
	DataStoreClient *datastore.Client
	Err             error
	Once            sync.Once
}

func getDataStoreKey(collection, id string) *datastore.Key {
	return datastore.NameKey(collection, id, nil)
}

func DefaultConfigFunc(projectId string, config *Configuration) {

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
