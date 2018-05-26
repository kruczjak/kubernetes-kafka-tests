#!/usr/bin/env bash

set -e

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker build -t kruczjak/kafka-example-app .
docker push kruczjak/kafka-example-app
rm main
