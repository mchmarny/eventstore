# Assumes following env vars set
#  GCP_PROJECT - ID of your project
#  CLUSTER_ZONE - GCP Zone, ideally same as your Knative k8s cluster
#  MYEVENTS_OAUTH_CLIENT_ID - Google OAuth2 Client ID
#  MYEVENTS_OAUTH_CLIENT_SECRET - Google OAuth2 Client Secret


# DEV
test:
	go test ./... -v

deps:
	go mod tidy

# BUILD
image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/myvents:latest

docker:
	docker build -t myvents .

# SERVICE
secrets:
	kubectl create secret generic myvents \
		--from-literal=OAUTH_CLIENT_ID=$(MYEVENTS_OAUTH_CLIENT_ID) \
		--from-literal=OAUTH_CLIENT_SECRET=$(MYEVENTS_OAUTH_CLIENT_SECRET)

service:
	kubectl apply -f deployments/service.yaml

service-clean:
	kubectl delete -f deployments/service.yaml
