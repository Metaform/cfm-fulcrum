# Development Deployment to Kind

## Docker Images

Build the CFM Fulcrum Agent Docker image.

## Kind Setup

Install Kind and create a cluster. Set it as the default K8S context. 

```
kind create cluster --config kind-config.yaml
```

Load the runtime images locally into Kind:

```
kind load docker-image pmanager:latest
kind load docker-image tmanager:latest
```
and
```
kind load docker-image cfm-fulcrum:latest
```

Note the Docker image will need to be reloaded each time it is updated.

## Terraform Deployment

Ensure the agent Docker image `pull_policy` is set to `Never` (the default). The Terraform scripts are
configured to use the default K8S context. To deploy:

```
terraform init
terraform apply
```