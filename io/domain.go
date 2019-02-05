package io

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/plancks-cloud/plancks-cloud/model"
)

const ServiceCollectionName = model.ServiceCollectionName

type Service model.Service

func ListAllServices(db *Configuration) (sl []*Service, err error) {
	ctx := context.Background()
	sl = make([]*Service, 0)
	q := datastore.NewQuery(ServiceCollectionName)
	_, err = db.DataStoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("DataStore DB: could not list services: %v", err)
	}
	return sl, nil
}
func StoreServices(db *Configuration, id string, sl []*Service) {
	key := getDataStoreKey(ServiceCollectionName, id)
	_, err := db.DataStoreClient.Put(context.Background(), key, sl)
	if err != nil {
		fmt.Println(err)
	}
}

const RouteCollectionName = model.RouteCollectionName

type Route model.Route

func ListAllRoutes(db *Configuration) (sl []*Route, err error) {
	ctx := context.Background()
	sl = make([]*Route, 0)
	q := datastore.NewQuery(RouteCollectionName)
	_, err = db.DataStoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("DataStore DB: could not list routes: %v", err)
	}
	return sl, nil
}
func StoreRoutes(db *Configuration, id string, sl []*Route) {
	key := getDataStoreKey(RouteCollectionName, id)
	_, err := db.DataStoreClient.Put(context.Background(), key, sl)
	if err != nil {
		fmt.Println(err)
	}
}
