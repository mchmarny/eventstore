# Assumes following env vars set
#  GCP_PROJECT - ID of your project
#  CLUSTER_ZONE - GCP Zone, ideally same as your Knative k8s cluster
#  MYEVENTS_OAUTH_CLIENT_ID - Google OAuth2 Client ID
#  MYEVENTS_OAUTH_CLIENT_SECRET - Google OAuth2 Client Secret


# DEV
test:
	go test ./... -v

build:
	go build ./... -v

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
		--from-literal=OAUTH_CLIENT_SECRET=${MYEVENTS_OAUTH_CLIENT_SECRET}

secrets-clean:
	kubectl delete secret myevents

service:
	kubectl apply -f deployments/service.yaml

service-clean:
	kubectl delete -f deployments/service.yaml
