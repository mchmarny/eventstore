package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	prj "github.com/mchmarny/gcputil/project"
	ev "github.com/mchmarny/gcputil/env"

)

const (
	defaultCollectionName = "cloudevents"
)

var (
	coll *firestore.CollectionRef
)

func init() {

	projectID := prj.GetIDOrFail()
	collName := ev.MustGetEnvVar("FIRESTORE_COLL_NAME", defaultCollectionName)

	log.Printf("Initiating firestore client for %s collection in %s project",
		collName, projectID)

	// Assumes GOOGLE_APPLICATION_CREDENTIALS is set
	dbClient, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Error while creating Firestore client: %v", err)
	}
	coll = dbClient.Collection(collName)
}

func saveData(ctx context.Context, id string, data interface{}) error {

	log.Printf("Saving id:%s - %v", id, data)

	if id == "" {
		log.Println("nil id on save")
		return errors.New("Nil ID")
	}

	_, err := coll.Doc(id).Set(ctx, data)
	if err != nil {
		log.Printf("error on save: %v", err)
		return fmt.Errorf("Error on save: %v", err)
	}

	return nil

}
