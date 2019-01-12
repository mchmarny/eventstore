package main

import (
	"context"
	"flag"
	"fmt"
	"bytes"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/url"
	"encoding/json"
    "net/http"

	"github.com/mchmarny/myevents/pkg/utils"
	"github.com/cloudevents/sdk-go/v02"
)

const (
	defaultPushEventingSource = "myevents-client"
	defaultNumberOfMessages   = 3
)

var (
	canceling     bool
	toke          string
	eventSrc      string
	targetURL     string
	numOfMessages int
)


func main() {

	// flags
	flag.StringVar(&toke, "toke", os.Getenv("MYEVENTS_KNOWN_PUBLISHER_TOKEN"), "Known publisher token")
	flag.StringVar(&eventSrc, "src", defaultPushEventingSource, "Source of data (Optional)")
	flag.StringVar(&targetURL, "url", "", "Target service URL where events will be sent")
	flag.IntVar(&numOfMessages, "messages", defaultNumberOfMessages, "Number of messages to sent [3]")
	flag.Parse()

	if toke == "" {
		log.Fatal("`token` required (or define MYEVENTS_KNOWN_PUBLISHER_TOKEN env var)")
	}

	if targetURL == "" {
		log.Fatal("`url` required ")
	}

	// context
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		log.Println(<-ch)
		canceling = true
		cancel()
		os.Exit(0)
	}()

	status := make(chan string)
	done := make(chan int)

	// start sending data
	go sendMessages(ctx, status, done)

F:
	for {
		select {
		case <-ctx.Done():
			break F
		case s := <-status:
			fmt.Printf("   status: %s \n", s)
		case c := <-done:
			fmt.Printf("sent %d messages \n", c)
			break F
		}
	}

}



func sendMessages(ctx context.Context, status chan<- string, done chan<- int) {

	srcURL, _ := url.Parse("https://github.com/mchmarny/myevents/cmd/client")
	sentCount := 0
	for i := 0; i < numOfMessages; i++ {

		now := time.Now().UTC()
		event := &v02.Event{
			SpecVersion: "0.2",
			Type:        "tech.knative.event.write",
			Source:      *srcURL,
			ID:          utils.MakeUUID(),
			Time: 		 &now,
			ContentType: "text/plain",
			Data: 		 fmt.Sprintf("message number %d", i),
		}

		data, _ := json.Marshal(event)
		err := postContent(event.ID, data)
		if err != nil {
			status <- fmt.Sprintf("error msg[%d] %s", i, err.Error())
		}
		sentCount++
	}

	done <- sentCount
	return

}


func postContent(id string, content []byte) error {
	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(content))
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
	log.Printf("post %s status: %s", id, resp.Status)
	return nil
}