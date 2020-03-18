package cmd

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func checkDB(mongoDBHost string, mongoDBPassword string) error {
	dbHost := strings.Replace(mongoDBHost, "<password>", mongoDBPassword, -1)
	fmt.Println(dbHost)
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
