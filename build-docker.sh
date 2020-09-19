#!/bin/bash

docker build -t custom-plugin .
docker run --name custom-plugin --rm -d -p 8000:8000 -p 8080:8080 custom-plugin