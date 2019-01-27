
# Assumes following env vars set
#  GCP_PROJECT - ID of your project
#  CLUSTER_ZONE - GCP Zone, ideally same as your Knative k8s cluster


.PHONY: app client service

# DEV
test:
	go test ./... -v

service:
	go build ./cmd/service/

run:
	go run ./cmd.service/*.go

deps:
	go mod tidy

# BUILD

image:
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/myevents:latest

# DEPLOYMENT

deployment:
	kubectl apply -f deployments/service.yaml

nodeployment:
	kubectl delete -f deployments/service.yaml

# DEMO

event:
	# TARGET_URL=https://events.default.knative.tech/
	curl -H "Content-Type: application/json" \
		 -X POST -d @test-event.json http://localhost:8080/ | jq '.'
