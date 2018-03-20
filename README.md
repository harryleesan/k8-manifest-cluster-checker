# Manifest K8 Checker Go

## Introduction

This Go app checks your local manifests (.yml files) against your Kubernetes
cluster and does a `diff` and outputs any discrepancies.

Currently, the files that will be checked **must be named** `deployment.yml` and `service.yml`
with **single** resource definition. Thus, only `Deployment` and `Service`
resources are checked.

## Requirements
Ensure that `go` and `kubectl` is installed on your system. Also make sure that
your `kubectl` is pointed at the correct cluster.

### Compile

```bash
go build setup.go main.go
```

This will generate a binary `main`.

### Usage

Place the compiled binary in the root directory of where you store your
Kubernetes manifests. If your manifests ares stored in `apps/`:

```bash
./main apps
```

The app is configured to **only** check `deployment.yml` and
`service.yml` where you would define `Deployment` and `Service` resources
respectively.

### TODO
- [ ] Reduce binary size (right now it's ~30MB)
- [ ] Give options to check for other resources
