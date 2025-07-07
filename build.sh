#!/bin/bash
rm arkan-face-key

CGO_ENABLED=1 \
CGO_CXXFLAGS="--std=c++14" \
CGO_LDFLAGS="-L/usr/local/lib -ljpeg -ldlib -lblas -lcblas -llapack" \
go build -o arkan-face-key

docker rm -f arkan-face-key
docker rmi arkan-face-key
docker build -t arkan-face-key .
docker run --name arkan-face-key --network dev --restart always -dp 9000:9000 -v /data/media/arkan/face_key/:/app/faces/images/ arkan-face-key