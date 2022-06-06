IMG ?= quay.io/olavtar/chromedp

DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
OUT_FILE := "$(DIR)chromedp"

build:
	CGO_ENABLED=0 go build -v

release: build docker-build docker-push

docker-build:
	docker build -t ${IMG} .

docker-push:
	docker push ${IMG}
