#!/bin/bash
podman run -d --name service-discovery -p 9098:9098 service-discovery
podman run -d --name dmsg-discovery -p 9090:9090 dmsg-discovery
podman run -d --name dmsg-server -p 8081:8081 dmsg-server
podman run -d --name address-resolver -p 9093:9093 address-resolver
podman run -d --name route-finder -p 9092:9092 route-finder
