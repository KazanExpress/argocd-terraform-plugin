BINARY=argocd-terraform-plugin

default: build

quality:
	go vet github.com/KazanExpress/argocd-terraform-plugin
	go test -v -coverprofile cover.out ./...

build:
	go build -o ${BINARY} .

install: build

e2e: install
	./argocd-terraform-plugin
