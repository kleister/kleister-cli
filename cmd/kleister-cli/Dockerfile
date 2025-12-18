FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine3.21@sha256:b4dbd292a0852331c89dfd64e84d16811f3e3aae4c73c13d026c4d200715aff6 AS builder

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

FROM alpine:3.23@sha256:be171b562d67532ea8b3c9d1fc0904288818bb36fc8359f954a7b7f1f9130fb2

RUN apk add --no-cache ca-certificates mailcap && \
    addgroup -g 1337 kleister && \
    adduser -D -u 1337 -h /var/lib/kleister -G kleister kleister

ENTRYPOINT ["/usr/bin/kleister-cli"]
CMD ["help"]

COPY --from=builder /go/src/kleister/bin/kleister-cli /usr/bin/kleister-cli
WORKDIR /var/lib/kleister
USER kleister
