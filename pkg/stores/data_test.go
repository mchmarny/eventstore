package stores

import (
	"testing"
)

func TestJobData(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestJobData")
	}

	InitDataStore()

	//ctx := context.Background()

	// err := saveJob(termReq)

	// req, err := SaveEvent(termReq.ID)

	// if err != nil {
	// 	t.Errorf("Error on job read: %v", err)
	// }

	// if req.ID != termReq.ID {
	// 	t.Errorf("Got invalid job: %v", req)
	// }

}
