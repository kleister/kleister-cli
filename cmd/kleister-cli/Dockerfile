FROM --platform=$BUILDPLATFORM golang:1.26.4-alpine@sha256:f23e8b227fb4493eabe03bede4d5a32d04092da71962f1fb79b5f7d1e6c2a17f AS builder

RUN apk add --no-cache -U git curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR /go/src/kleister
COPY . /go/src/kleister/

RUN --mount=type=cache,target=/go/pkg \
    go mod download -x

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    task generate build GOOS=${TARGETOS} GOARCH=${TARGETARCH}

FROM alpine:3.24@sha256:660e0827bd401543d81323d4886abbd08fda0fe3ba84337837d0b11a67251283

RUN apk add --no-cache ca-certificates mailcap && \
    addgroup -g 1337 kleister && \
    adduser -D -u 1337 -h /var/lib/kleister -G kleister kleister

ENTRYPOINT ["/usr/bin/kleister-cli"]
CMD ["help"]

COPY --from=builder /go/src/kleister/bin/kleister-cli /usr/bin/kleister-cli
WORKDIR /var/lib/kleister
USER kleister
