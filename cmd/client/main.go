package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/mchmarny/myevents/pkg/clients"
)

const (
	defaultPushEventingSource = "myevents-client"
	defaultNumberOfMessages   = 3
)

var (
	canceling bool
	toke      string
	eventSrc  string
	targetURL string
	message   string
	sender    *clients.Sender
)

func main() {

	// flags
	flag.StringVar(&toke, "toke", os.Getenv("MYEVENTS_KNOWN_PUBLISHER_TOKEN"), "Known publisher token")
	flag.StringVar(&targetURL, "url", "", "Target service URL where events will be sent")
	flag.StringVar(&eventSrc, "src", defaultPushEventingSource, "Source of data (Optional)")
	flag.StringVar(&message, "message", "test message", "The content of the message [test message]")
	flag.Parse()

	if toke == "" {
		log.Fatal("`token` required (or define MYEVENTS_KNOWN_PUBLISHER_TOKEN env var)")
	}

	if targetURL == "" {
		log.Fatal("`url` required ")
	}

	if message == "" {
		log.Fatal("`message` required ")
	}

	// context
	ctx := context.Background()

	var err error
	sender, err = clients.NewSender(targetURL)
	if err != nil {
		log.Fatalf("error while creating sender: %v", err)
	}

	if err = sender.SendMessages(ctx, "tech.knative.event.write", message); err != nil {
		log.Fatalf("error while sending: %v", err)
	}

}
