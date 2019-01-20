
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
	./service

deps:
	go mod tidy

# BUILD

image:
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/myevents:latest

docker:
	docker build -t myevents .

# DEPLOYMENT

deployment:
	kubectl apply -f deployments/service.yaml

nodeployment:
	kubectl delete -f deployments/service.yaml

# DEMO

event:
	curl -H "Content-Type: application/json" \
		 -X POST --data "{ \
			\"specversion\": \"0.2\", \
			\"type\": \"tech.knative.event.write\", \
			\"source\": \"https://knative.tech/test\", \
			\"id\": \"id-0000-1111-2222-3333-4444\", \
			\"time\": \"2019-01-11T17:31:00Z\", \
			\"contenttype\": \"text/plain\", \
			\"data\": \"My message content\" \
		}" \
		"http://localhost:8080/" \
		| jq '.'

client:
	go build ./cmd/client/

client-event:
	./client --message "Test message" --url "http://events.default.knative.tech/" | jq '.'

client-event-local:
	./client --message "Test message" --url "http://localhost:8080/" | jq '.'