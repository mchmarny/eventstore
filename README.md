# eventstore

Simple Knative service persisting Cloud Events to Cloud Firestore collection. Useful in Knative Events demos

## Prerequisites

 * [Knative](https://github.com/knative/docs/blob/master/install) installed
    * Configured [outbound network access (https://github.com/knative/docs/blob/master/serving/outbound-network-access.md)
    * Installed [Knative Eventing](https://github.com/knative/docs/tree/master/eventing) using the `release.yaml` file


## Deployment

Firestore client still requires GCP Project ID to create a client. So, before we can deploy this service to Knative, you will need to update the `GCP_PROJECT_ID` in Now in the `service.yaml` file.

```yaml
    - name: GCP_PROJECT_ID
      value: "enter your project ID here"
```

Once done updating our service manifest (`service.yaml`) you are ready to deploy it.

```shell
kubectl apply -f deployments/service.yaml -n demo
```

The response should be

```shell
service.serving.knative.dev "eventstore" configured
```

To check if the service was deployed successfully you can check the status using `kubectl get pods -n demo` command. The response should look something like this (e.g. Ready `3/3` and Status `Running`).

```shell
NAME                                          READY     STATUS    RESTARTS   AGE
eventstore-0000n-deployment-5645f48b4d-mb24j  3/3       Running   0          10s
```

## Configuration

To make `eventstore` service `cluster.local` so it does not expose externally accessible endpoint but still enable other services to discover it using `eventstore` simply add `serving.knative.dev/visibility: cluster-local` label to `deployments/service.yaml` service manifest

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

