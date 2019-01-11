# Assumes following env vars set
#  GCP_PROJECT - ID of your project
#  CLUSTER_ZONE - GCP Zone, ideally same as your Knative k8s cluster
#  GAUTHER_OAUTH_CLIENT_ID - Google OAuth2 Client ID
#  GAUTHER_OAUTH_CLIENT_SECRET - Google OAuth2 Client Secret


# DEV
test:
	go test ./... -v

deps:
	go mod tidy


# REDIS
redis-secret:
	kubectl create secret generic env-secrets --from-literal=REDIS_PASS=$(REDIS_PASS)

redis-disk:
	gcloud compute --project=$(GCP_PROJECT) disks create \
		redis-disk --zone=$(CLUSTER_ZONE) --type=pd-ssd --size=10GB

redis:
	kubectl apply -f deployments/redis-pd.yaml

# BUILD
image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/kueue:latest

docker:
	docker build -t kueue .

# SERVICE
secrets:
	kubectl create secret generic kueue \
		--from-literal=OAUTH_CLIENT_ID=$(KUEUE_OAUTH_CLIENT_ID) \
		--from-literal=OAUTH_CLIENT_SECRET=$(KUEUE_OAUTH_CLIENT_SECRET)

service:
	kubectl apply -f deployments/service.yaml
	kubectl get pods

service-clean:
	kubectl delete -f deployments/service.yaml
	kubectl get pods
