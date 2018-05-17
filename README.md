# Kubernetes Manifest Checker

## Introduction

This **Go** app checks your local manifests (`.yml` files) against your Kubernetes
cluster and does a `diff` to display any discrepancies. This only works for
manifests that have been created using `kubectl apply`.

Currently, the files that will be checked **must be named** `deployment.yml` and `service.yml`
with a **single** resource definition in each file. Thus, only `Deployment` and `Service`
resources are checked.

## Requirements
Ensure that `go`, `glide` and `kubectl` is installed on your system. Also make sure that
your `kubectl` is pointed at the correct cluster.

### Compile
Ensure that this project is in your `$GOPATH/src/`.

```bash
glide install
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
respectively. Other resource types will be added in the future.

### TODO
- [ ] Reduce binary size (right now it's ~30MB)
- [ ] Give options to check for other resources
