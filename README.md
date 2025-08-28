# Echo Operator

A Kubernetes operator built with Go and operator-sdk that demonstrates operator development patterns.

## Overview

The Echo Operator takes a text message and creates a Kubernetes Job that runs a Pod to echo the specified text. This project serves as a practical learning example for:

- Go programming fundamentals
- Kubernetes operator development with operator-sdk
- Custom Resource Definitions (CRDs)
- Controller patterns and reconciliation loops
- Kubernetes Job management

## Features

- **Custom Resource**: `Echo` CRD with message validation
- **Job Management**: Creates and monitors Kubernetes Jobs
- **Status Tracking**: Phase-based status with Kubernetes conditions
- **Validation**: Input validation with proper error messages
- **Clean CLI**: Custom print columns for `kubectl get echo`

## Quick Start

```bash
# Apply the CRD and RBAC
make install

# Run the operator locally
make run

# In another terminal, create an Echo resource
kubectl apply -f config/samples/echo_v1alpha1_echo.yaml

# Check the status
kubectl get echo
kubectl describe echo echo-sample
```

## Development

### Local Development

```bash
# Run tests
make test

# Generate code and manifests
make generate manifests

# Build binary
make build

# Run locally (requires kubeconfig)
make run
```

## Architecture
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Echo CRD      │───▶│ Echo Controller  │───▶│ Kubernetes Job  │
│                 │    │                  │    │                 │
│ spec:           │    │ Reconcile Loop   │    │ Runs Pod with   │
│   message: "hi" │    │                  │    │ echo command    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

The controller watches for Echo resources and creates Jobs to execute the echo command with the specified message.
