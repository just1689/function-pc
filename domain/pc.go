package domain

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/plancks-cloud/function-pc/io"
	"github.com/plancks-cloud/plancks-cloud/model"
	"github.com/sirupsen/logrus"
)

const ServiceCollectionName = model.ServiceCollectionName

type Service model.Service

func ListAllServices(db *io.Configuration) (sl []Service, err error) {
	ctx := context.Background()
	sl = make([]Service, 0)
	q := datastore.NewQuery(ServiceCollectionName)
	_, err = db.DataStoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("DataStore DB: could not list services: %v", err)
	}
	return sl, nil
}
func StoreServices(db *io.Configuration, id string, sl []Service) (err error) {
	key := io.GetDataStoreKey(ServiceCollectionName, id)
	_, err = db.DataStoreClient.Put(context.Background(), key, sl)
	if err != nil {
		logrus.Println(err)
	}
	return
}

const RouteCollectionName = model.RouteCollectionName

type Route model.Route

func ListAllRoutes(db *io.Configuration) (sl []Route, err error) {
	ctx := context.Background()
	sl = make([]Route, 0)
	q := datastore.NewQuery(RouteCollectionName)
	_, err = db.DataStoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("DataStore DB: could not list routes: %v", err)
	}
	return sl, nil
}
func StoreRoutes(db *io.Configuration, id string, sl []Route) (err error) {
	key := io.GetDataStoreKey(RouteCollectionName, id)
	_, err = db.DataStoreClient.Put(context.Background(), key, sl)
	if err != nil {
		logrus.Println(err)
	}
	return
}