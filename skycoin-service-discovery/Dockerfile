FROM golang:1.19-alpine AS builder

ARG BUILDINFO_LDFLAGS
ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED} \
    GOOS=linux  \
    GO111MODULE=on

COPY . /service-discovery

WORKDIR /service-discovery

# Build the discovery
RUN apk add --no-cache make git && \
    sh -c /service-discovery/ci_scripts/install-preq.sh && \
    make build-deploy

# Build image
FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add --no-cache ca-certificates openssl iproute2 git
RUN update-ca-certificates --fresh

COPY --from=builder /release/service-discovery /usr/local/bin/service-discovery
ENTRYPOINT ["service-discovery"]
