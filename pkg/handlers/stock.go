package handlers

import (
	"log"
	"context"
	"errors"

	"github.com/mchmarny/myevents/pkg/stores"
	"github.com/mchmarny/kapi/common"
	"github.com/knative/pkg/cloudevents"

)

const (
	knownPublisherTokenName = "token"
)

// StockHandler submitted messages
func StockHandler(ctx context.Context, e *common.SimpleStock) error {

	ec := cloudevents.FromContext(ctx)
	if ec != nil {
		log.Printf("Received Cloud Event Context as: %+v", *ec)
	} else {
		log.Printf("No Cloud Event Context found")
	}

	log.Printf("Stock %v", e)

	if e.ID == "" {
		log.Println("unable to parse event ID")
		return errors.New("Invalid event format")
	}


	saveErr := stores.SaveEvent(ctx, e.ID, e)
	if saveErr != nil {
		log.Printf("error on event save: %v", saveErr)
		return errors.New( "Error on event save")
	}

	log.Printf("Event saved: %v", e.ID)
	return nil

}
