package main

import (
	"context"
	"flag"

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
	flag.StringVar(&targetURL, "url", "", "Target service URL where events will be sent")
	flag.StringVar(&eventSrc, "src", defaultPushEventingSource, "Source of data (Optional)")
	flag.StringVar(&message, "message", "test message", "The content of the message [test message]")
	flag.Parse()

	if targetURL == "" {
		panic("`url` required")
	}

	if message == "" {
		panic("`message` required")
	}

	// context
	ctx := context.Background()

	var err error
	sender, err = clients.NewSender(targetURL)
	if err != nil {
		panic("error while creating sender: " + err.Error())
	}

	if err = sender.SendMessages(ctx, "tech.knative.events.test", message); err != nil {
		panic("error while sending: " + err.Error())
	}

}
