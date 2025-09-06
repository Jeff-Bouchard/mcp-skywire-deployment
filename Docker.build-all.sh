#!/bin/bash

podman build -t service-discovery -f Dockerfile.service-discovery .
podman build -t dmsg-discovery -f Dockerfile.dmsg-discovery .
podman build -t dmsg-server -f Dockerfile.dmsg-server .
podman build -t address-resolver -f Dockerfile.address-resolver .
podman build -t route-finder -f Dockerfile.route-finder .
