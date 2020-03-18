package cmd

import (
	"bytes"
	"context"
	"io"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func checkFirebaseStorageBucket(pathWD string, bucketName string) error {
	ctx := context.Background()
	config := &firebase.Config{
		StorageBucket: bucketName + ".appspot.com",
	}
	opt := option.WithCredentialsFile(pathWD + "/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return err
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return err
	}
	bucket, err := client.DefaultBucket()
	if err != nil {
		return err
	}
	r := bytes.NewReader([]byte("Test firebase."))
	wc := bucket.Object("telar.test").NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func checkFirebaseServiceAccount(pathWD string) error {
	var file, err = os.OpenFile(pathWD+"/serviceAccountKey.json", os.O_RDONLY, 0444)
	if isError(err) {
		return err
	}
	defer file.Close()
	return nil
}
