package controller

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/plancks-cloud/function-pc/domain"
	"github.com/plancks-cloud/function-pc/io"
	"github.com/sirupsen/logrus"
)

func ListAllServices(db *io.Configuration) (sl []domain.Service, err error) {
	ctx := context.Background()
	sl = make([]domain.Service, 0)
	q := datastore.NewQuery(domain.ServiceCollectionName)
	_, err = db.DataStoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("DataStore DB: could not list services: %v", err)
	}
	return sl, nil
}
func StoreServices(db *io.Configuration, id string, sl []domain.Service) (err error) {
	key := io.GetDataStoreKey(domain.ServiceCollectionName, id)
	_, err = db.DataStoreClient.Put(context.Background(), key, &sl)
	if err != nil {
		logrus.Println(err)
	}
	return
}

func ListAllRoutes(db *io.Configuration) (sl []domain.Route, err error) {
	ctx := context.Background()
	sl = make([]domain.Route, 0)
	q := datastore.NewQuery(domain.RouteCollectionName)
	_, err = db.DataStoreClient.GetAll(ctx, q, &sl)
	if err != nil {
		return nil, fmt.Errorf("DataStore DB: could not list routes: %v", err)
	}
	return sl, nil
}
func StoreRoutes(db *io.Configuration, id string, sl []domain.Route) (err error) {
	key := io.GetDataStoreKey(domain.RouteCollectionName, id)
	_, err = db.DataStoreClient.Put(context.Background(), key, &sl)
	if err != nil {
		logrus.Println(err)
	}
	return
}
