#!/bin/bash

docker build -t minesweeperx .
docker run --rm -p 3000:3000 localhost/minesweeperx:latest
