package stores

import (
	"errors"
	"log"
	"context"
	"fmt"
	"cloud.google.com/go/firestore"

	"github.com/mchmarny/gauther/utils"


)

const (
	defaultCollectionName = "events"
)

var (
	coll   *firestore.CollectionRef
)



// InitDataStore initializes client
func InitDataStore() {

	projectID := utils.MustGetEnv("GCP_PROJECT_ID", "")
	collName := utils.MustGetEnv("FIRESTORE_COLL_NAME", defaultCollectionName)

	log.Printf("Initiating firestore client for %s collection in %s project",
		collName, projectID)

	// Assumes GOOGLE_APPLICATION_CREDENTIALS is set
	dbClient, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Error while creating Firestore client: %v", err)
	}
	coll = dbClient.Collection(collName)
}


// SaveEvent saves passed event
func SaveEvent(ctx context.Context, id string, data map[string]interface{}) error {

	if id == "" {
		return errors.New("Nil event ID")
	}

	_, err := coll.Doc(id).Set(ctx, data)
	if err != nil {
		return fmt.Errorf("Error on save: %v", err)
	}

	return nil

}
