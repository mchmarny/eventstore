# Assumes following env vars set
#  GCP_PROJECT - ID of your project

.PHONY: run mod image service event

# DEV

run:
	go run *.go -v

# BUILD

mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/myevents:0.1.3

sample-image: mod
	gcloud builds submit \
		--project knative-samples \
		--tag gcr.io/knative-samples/myevents:0.1.3

# DEPLOYMENT

service:
	kubectl apply -f service.yaml -n demo

# DEMO

event:
	curl -H "Content-Type: application/json" \
		 -H "CE-Specversion: 0.2" \
		 -H "CE-ID: 1111-2222-3333-4444-5555-6666" \
		 -H "CE-Type: com.twitter" \
		 -H "CE-Time: 2018-04-05T03:56:24Z" \
		 -H "CE-Source: https://twitter.com/api/1" \
		 -X POST -d @test-event.json http://localhost:8080/

