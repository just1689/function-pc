package io

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"sync"
)

type Configuration struct {
	DataStoreClient *datastore.Client
	Err             error
	Once            sync.Once
}

func (db *Configuration) ListAllServices() (sl []*Service, err error) {
	ctx := context.Background()
	sl = make([]*Service, 0)
	q := datastore.NewQuery(ServiceCollectionName)
	_, err = db.DataStoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("DataStore DB: could not list services: %v", err)
	}
	return sl, nil
}

func (db *Configuration) ListAllRoutes() (sl []*Route, err error) {
	ctx := context.Background()
	sl = make([]*Route, 0)
	q := datastore.NewQuery(RouteCollectionName)
	_, err = db.DataStoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("DataStore DB: could not list routes: %v", err)
	}
	return sl, nil
}

func (db *Configuration) StoreObjects(collection string, clones []*interface{}) {
	keys := getDataStoreKeys(collection, clones)
	_, err := config.DataStoreClient.PutMulti(context.Background(), keys, &clones)
	if err != nil {
		fmt.Println(err)
	}
}

func getDataStoreKeys(project string, clones *io.Clones) (keys []*datastore.Key) {
	keys = make([]*datastore.Key, len(clones.Days))
	for _, c := range clones.Days {
		key := datastore.NameKey(project, c.Timestamp, nil)
		keys = append(keys, key)
	}
	return keys
}
