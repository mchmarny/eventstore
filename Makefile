
# Assumes following env vars set
#  GCP_PROJECT - ID of your project
#  CLUSTER_ZONE - GCP Zone, ideally same as your Knative k8s cluster
#  MYEVENTS_OAUTH_CLIENT_ID - Google OAuth2 Client ID
#  MYEVENTS_OAUTH_CLIENT_SECRET - Google OAuth2 Client Secret
#  MYEVENTS_KNOWN_PUBLISHER_TOKEN - One of the known publisher tokens


.PHONY: app client service

# DEV
test:
	go test ./... -v

setup:
	export OAUTH_CLIENT_ID=$MYEVENTS_OAUTH_CLIENT_ID
	export OAUTH_CLIENT_SECRET=$MYEVENTS_OAUTH_CLIENT_SECRET
	export KNOWN_PUBLISHER_TOKENS=$MYEVENTS_KNOWN_PUBLISHER_TOKEN

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

secrets:
	kubectl create secret generic myevents \
		--from-literal=OAUTH_CLIENT_ID=${MYEVENTS_OAUTH_CLIENT_ID} \
		--from-literal=OAUTH_CLIENT_SECRET=${MYEVENTS_OAUTH_CLIENT_SECRET} \
		--from-literal=KNOWN_PUBLISHER_TOKENS=${KNOWN_PUBLISHER_TOKENS}

secrets-clean:
	kubectl delete secret myevents

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
		"http://localhost:8080/?token=${MYEVENTS_KNOWN_PUBLISHER_TOKEN}" \
		| jq '.'

client:
	go build ./cmd/client/

client-send:
	./client \
		--message "Test message from ${HOSTNAME}" \
		--url "https://myevents.default.knative.tech/?token=${MYEVENTS_KNOWN_PUBLISHER_TOKEN}"