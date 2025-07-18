# Running the Demo

The demo illustrates how CFM integrates with the Fulcrum Core service using a Fulcrum agent that processes `jobs`
created from Fulcrum `Service` requests.

## Steps

### Deploy the Services

Deploy the Fulcrum and CFM services to a Kubernetes cluster using Kind by following
the [instructions](../../deployments/terraform/environments/kind).

### Run the Setup

Fulcrum and CFM need to be seeded with test data. From the root directory, run the following commandline tool:

```
go run cmd/demo/main.go -action=onboard
```

## Run Demo Scenario

View the CFM Test Agent logs:

```
k logs -f -l app=testagent
```

From the root directory, run the following commandline tool to create a Fulcrum service:

```
go run cmd/demo/main.go -action=service
```

Fulcrum Core will queue a`Job` that will be picked up and processed by the CFM Fulcrum Agent. The agent will create a
CFM `Deployment` that will be sent to the CFM Activity Agent (`testagent`) for processing. 