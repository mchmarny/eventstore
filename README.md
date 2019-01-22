# myevents

Knative Cloud Event (v0.2) Collector with Cloud Firestore persistence


## Prerequisites

 * `gcloud` configured. If not, see [Installing Google Cloud SDK](https://cloud.google.com/sdk/install)
 * [Knative](https://github.com/knative/docs/blob/master/install) installed
    * Configured [outbound network access (https://github.com/knative/docs/blob/master/serving/outbound-network-access.md)
    * Installed [Knative Eventing](https://github.com/knative/docs/tree/master/eventing) using the `release.yaml` file


> This readme is still a bit of work in progress so if you are finding something missing do take a look at the [Makefile](https://github.com/mchmarny/myevents/blob/master/Makefile)

## Deployment

To deploy the `myevents` are are going to:

- [myevents](#myevents)
  - [Prerequisites](#prerequisites)
  - [Deployment](#deployment)
    - [Build the image](#build-the-image)
    - [Configure Knative](#configure-knative)
    - [Deploy Service](#deploy-service)
  - [Disclaimer](#disclaimer)

### Build the image

Quickest way to build your service image is through [GCP Build](https://cloud.google.com/cloud-build/). Just submit the build request from within the `myevents` directory:

```shell
gcloud builds submit \
    --project ${GCP_PROJECT} \
	--tag gcr.io/${GCP_PROJECT}/myevents:latest
```

The build service is pretty verbose in output but eventually you should see something like this

```shell
ID           CREATE_TIME          DURATION  SOURCE                                   IMAGES                      STATUS
6905dd3a...  2018-12-23T03:48...  1M43S     gs://PROJECT_cloudbuild/source/15...tgz  gcr.io/PROJECT/myevents SUCCESS
```

Copy the image URI from `IMAGE` column (e.g. `gcr.io/PROJECT/myevents`).

### Configure Knative

Before we can deploy that service to Knative, we just need to update the `GCP_PROJECT_ID` in Now in the `deployments/service.yaml` file.

```yaml
    - name: GCP_PROJECT_ID
      value: "enter your project ID here"
```

> Note, if you want to be able to post to this service from external clients, remove the `serving.knative.dev/visibility: cluster-local` label in `deployments/service.yaml`


### Deploy Service

Once done updating our service manifest (`deployments/service.yaml`) you are ready to deploy it.

```shell
kubectl apply -f deployments/service.yaml
```

The response should be

```shell
service.serving.knative.dev "myevents" configured
```

To check if the service was deployed successfully you can check the status using `kubectl get pods` command. The response should look something like this (e.g. Ready `3/3` and Status `Running`).

```shell
NAME                                          READY     STATUS    RESTARTS   AGE
myevents-0000n-deployment-5645f48b4d-mb24j     3/3       Running   0          4h
```

> Note, the `myevents` service is cluster.local only and does not expose externally accessible endpoint. Other service can publish to it using its service name `events`.

To make the service externally accessible remove the `serving.knative.dev/visibility: cluster-local` label from `deployments/service.yaml` service manifest

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

