SHELL := /bin/bash

.PHONY: semantick

run: 
	./semantick

build:
	go build -ldflags "-X main.build=local"

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

all: semantick 

semantick:
	docker build \
		-f zarf/docker/Dockerfile.semantick \
		-t semantick-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := semantick-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current 


kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	kind load docker-image semantick-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	cat zarf/k8s/base/semantick-pod/base-service.yaml | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-logs:
	kubectl logs -l app=semantick --all-containers=true -f --tail=100  

# TODO: Get rid of the --namespace tag by setting the kubectl config in a 
# previous step.
kind-restart:
	kubectl rollout restart deployment semantick-pod --namespace=semantick-system

kind-update: all kind-load kind-restart

kind-describe:
	kubectl describe pod -l app=semantick
