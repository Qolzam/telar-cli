package cmd

import (
	"context"

	"github.com/Qolzam/telar-cli/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func checkDB(mongoDBURI string) error {
	dbHost := mongoDBURI
	log.Info(dbHost)
	client, err := mongo.NewClient(options.Client().ApplyURI(dbHost))
	ctx := context.Background()
	err = client.Connect(ctx)
	if isError(err) {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if isError(err) {
		return err
	}
	return nil
}
