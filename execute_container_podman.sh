#!/bin/bash

podman build -t minesweeperx .
podman run --rm -p 3000:3000 localhost/minesweeperx:latest
