LDFLAGS += -X main.version=$$(git describe --always --abbrev=40 --dirty)
TEST?=$$(go list ./... |grep -v 'vendor')
PKG_NAME=openshift
TERRAFORM_PLUGINS=$(HOME)/.terraform.d/plugins

default: fmt build

build:
	go build -ldflags "${LDFLAGS}" -tags "${TAGS}"

install: default
	mkdir -p ${TERRAFORM_PLUGINS}
	mv terraform-provider-openshift ${TERRAFORM_PLUGINS}

fmt:
	go fmt ./openshift .

tools:
	go get golang.org/x/lint/golint

lint: tools
	go run golang.org/x/lint/golint -set_exit_status ./openshift .

test:
	go test -tags "${TAGS}" -v ./openshift

clean:
	rm -f terraform-provider-openshift

.PHONY: build install test fmt lint
