package clients

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	ce "github.com/knative/pkg/cloudevents"
	"github.com/mchmarny/myevents/pkg/utils"
)

const (
	sourceURI = "https://github.com/mchmarny/myevents/client"
)

// NewSender creates a already preconfigured Sender
func NewSender(targerURL string) (sender *Sender, err error) {

	s := &Sender{
		TargerURL: targerURL,
		SourceURL: sourceURI,
	}

	return s, nil

}

// Sender sends messages
type Sender struct {
	TargerURL string
	SourceURL string
}

// SendMessages sends v02.Event based on the provided data
func (s *Sender) SendMessages(ctx context.Context, eventType, text string) error {

	ex := ce.EventContext{
		CloudEventsVersion: "0.2",
		EventID:            utils.MakeUUID(),
		EventTime:          time.Now().UTC(),
		EventType:          eventType,
		EventTypeVersion:   "v0.1",
		ContentType:        "text/plain",
		Source:             s.SourceURL,
	}

	req, err := ce.Binary.NewRequest(s.TargerURL, text, ex)
	if err != nil {
		log.Printf("Error creating new quest: %v", err)
		return err
	}

	log.Printf("Posting to %s: %v", s.TargerURL, text)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted {

		log.Printf("Response Status: %s", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Response Body: %s", string(body))

		return fmt.Errorf("Send returned an invalid status: %s", resp.Status)
	}

	return nil

}
