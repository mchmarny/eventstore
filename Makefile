# Assumes following env vars set
#  GCP_PROJECT - ID of your project
#  CLUSTER_ZONE - GCP Zone, ideally same as your Knative k8s cluster
#  MYEVENTS_OAUTH_CLIENT_ID - Google OAuth2 Client ID
#  MYEVENTS_OAUTH_CLIENT_SECRET - Google OAuth2 Client Secret


# DEV
test:
	go test ./... -v

app:
	rm ./app
	go build ./cmd/app/
	./app

deps:
	go mod tidy

# BUILD
image:
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/myevents:latest

docker:
	docker build -t myevents .

# REDIS
redis-secret:
	kubectl create secret generic env-secrets --from-literal=REDIS_PASS=${REDIS_PASS}

redis-disk:
	gcloud compute --project=${GCP_PROJECT} disks create \
		redis-disk --zone=${CLUSTER_ZONE} --type=pd-ssd --size=10GB

redis:
	kubectl apply -f deployments/redis-pd.yaml

# SERVICE
secrets:
	kubectl create secret generic myevents \
		--from-literal=OAUTH_CLIENT_ID=${MYEVENTS_OAUTH_CLIENT_ID} \
		--from-literal=OAUTH_CLIENT_SECRET=${MYEVENTS_OAUTH_CLIENT_SECRET} \
		--from-literal=KNOWN_PUBLISHER_TOKENS=${KNOWN_PUBLISHER_TOKENS}

secrets-clean:
	kubectl delete secret myevents

service:
	kubectl apply -f deployments/service.yaml

service-clean:
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
		"https://myevents.default.knative.tech/v1/event?token=${MYEVENTS_KNOWN_PUBLISHER_TOKEN}" \
		| jq '.'

client:
	rm ./client
	go build ./cmd/client/
	./client --url "https://myevents.default.knative.tech/v1/event?token=${MYEVENTS_KNOWN_PUBLISHER_TOKEN}"