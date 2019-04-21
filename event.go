package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	ce "github.com/cloudevents/sdk-go"
)

const (
	idPrefix = "eid-"
)

type eventData struct {
	Context interface{}            `firestore:"ctx"`
	Data    map[string]interface{} `firestore:"data"`
}

type eventReceiver struct{}

func (r *eventReceiver) Receive(ctx context.Context, event ce.Event, _ *ce.EventResponse) error {

	log.Printf("Raw Event: %v", event)
	// event.DataAs(&data); err != nil {

	if event.ID() == "" {
		log.Println("unable to parse event ID")
		return errors.New("Invalid event format")
	}

	// Workaround for Firestore ID foramt requirements
	// can't stat with number
	eid := fmt.Sprintf("%s%s", idPrefix, event.ID())

	var p map[string]interface{}
	if err := event.DataAs(&p); err != nil {
		log.Printf("Failed to DataAs: %s", err.Error())
		return err
	}

	d := &eventData{Context: event.Context.AsV02(), Data: p}

	return saveData(ctx, eid, d)

}
