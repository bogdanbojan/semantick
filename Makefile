SHELL := /bin/bash

.PHONY: semantick

# ==============================================================================
# Testing running system

# Access metrics directly (4000)
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# For testing load on the service.
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem

# Testing Authorization
# curl -il https://localhost:3000/v1/testauth
# curl -il -H "Authorization: Bearer ${TOKEN}" https://localhost:3000/v1/testauth

# Database Access
# dblab --host 0.0.0.0 --user postgres --db postgres --pass postgres --ssl disable --port 5432 --driver postgres
# ==============================================================================

run: 
	go run app/services/semantick/main.go | go run app/tooling/logfmt/main.go

build:
	go build -ldflags "-X main.build=local"

# ==============================================================================
# Building containers

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
	kubectl config set-context --current --namespace=semantick-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/semantick-pod; kustomize edit set image semantick-image=semantick-amd64:$(VERSION)
	kind load docker-image semantick-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/semantick-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-semantick:
	kubectl get pods -o wide --watch --namespace=semantick-system

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-logs:
	kubectl logs -l app=semantick --all-containers=true -f --tail=100  | go run app/tooling/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment semantick-pod 

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=semantick


# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck -checks=all ./...

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all
