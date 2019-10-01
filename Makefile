FAAS_GATEWAY ?= http://127.0.0.1:8080

.PHONY: build all dev

all: build
	faas-cli up --build-arg GO111MODULE=on -f break.yml -g $(FAAS_GATEWAY)

build: template
	faas-cli build --build-arg GO111MODULE=on -f break.yml

template:
	faas-cli template pull https://github.com/fopina/golang-http-template.git --overwrite

dev: export BUILD_ENV=-dev
dev: export ZEROSCALE=true
dev: all

localweb: build
	docker run -p 9999:8082 --rm -ti fopina/functions:break
