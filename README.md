# Kubernetes Manifest Checker

## Introduction

This **Go** app checks your local manifests (`.yml` files) against your Kubernetes
cluster and does a `diff` to display any discrepancies. This only works for
manifests that have been created using `kubectl apply`.

Currently, the files that will be checked **must contain a single** resource definition in each file.

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

The app is configured to **only** check the resources that are specified.
Currently, these resources are `deployments`, `services` and `ingresses`.

### TODO
- [ ] Reduce binary size (right now it's ~30MB).
- [ ] Refactor code to allow other resources.
- [ ] Refactor code to allow for multiple resource definitions in the same YAML
  file.
