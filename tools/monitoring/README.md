# Knative Monitoring

Knative monitoring tool checks the prow job status and scrape through all the failure logs to catch test infrastructure
failures.

## Setup

### Create the cluster

```bash
gcloud container clusters create monitoring --enable-ip-alias
```

## Build and Deploy Changes

* `images/monitoring/Makefile` Commands to build and deploy the monitoring images

* `tools/monitoring/gke_deployment` YAML configuration to setup all the Kubernetes resources. Use `kubectl apply` to apply the changes.
