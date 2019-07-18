#!/bin/bash

cd $(dirname $0)
cd ..

IMG=registry.gitlab.com/fopina/caixabreak
VERSION=$(git rev-parse --short HEAD)

docker build -t $IMG -f docker/Dockerfile .
docker tag $IMG $IMG:$VERSION

if [ -z "$1" ]; then
	docker push $IMG:latest
	docker push $IMG:$VERSION
fi