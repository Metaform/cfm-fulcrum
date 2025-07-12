# Fulcrum CFM Agent

A [Fulcrum](https://github.com/fulcrumproject/core) agent
for [CFM](https://github.com/Metaform/connector-fabric-manager).

## Development Setup

### Requirements

Go 1.24.4, a Docker client, and Terraform

### Workspace

This project requires the following modules from CFM:

```
https://github.com/Metaform/connector-fabric-manager/common
https://github.com/Metaform/connector-fabric-manager/assembly
```

You may want to create a Go Workspace and have dependency resolution redirect to a local copy for development.

### Deployment

The system can be deployed to a Kind cluster using Terraform by following the steps
outlined [here.](/deployments/terraform/environments/kind/README.md)

