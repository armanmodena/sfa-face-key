#!/bin/bash
rm arkan-face-key-dev

CGO_ENABLED=1 \
CGO_CXXFLAGS="--std=c++14" \
CGO_LDFLAGS="-L/usr/local/lib -ljpeg -ldlib -lblas -lcblas -llapack" \
go build -o arkan-face-key-dev

docker rmi altercode99/arkan-face-key-dev:1.0.5
docker build -f Dockerfile.dev -t altercode99/arkan-face-key-dev:1.0.6 .
