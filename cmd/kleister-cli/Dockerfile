FROM --platform=$BUILDPLATFORM debian:13-slim@sha256:8aa5a7bf5a1955bddbe7d74fcdd0b91ab30e74362a8b73af672b0481b35b44a0 AS builder
WORKDIR /src

RUN apt-get update && \
    apt-get -y --no-install-recommends install ca-certificates curl git build-essential && \
    rm -rf /var/lib/apt/lists/*

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV PATH="/root/.local/bin:/mise/shims:$PATH" \
    MISE_DATA_DIR="/mise" \
    MISE_CONFIG_DIR="/mise" \
    MISE_CACHE_DIR="/mise/cache" \
    MISE_INSTALL_PATH="/usr/local/bin/mise"

RUN curl --proto '=https' --proto-redir '=https' \
    --fail --show-error --silent --location https://mise.run | sh

COPY mise.toml ./
RUN --mount=type=cache,target=/mise/cache \
   mise trust && mise install

COPY . /src/
RUN --mount=type=cache,target=/root/go/pkg \
    go mod download -x

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    task generate build GOOS=${TARGETOS} GOARCH=${TARGETARCH}

FROM alpine:3.24@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf

RUN apk add --no-cache ca-certificates mailcap && \
    addgroup -g 1337 kleister && \
    adduser -D -u 1337 -h /var/lib/kleister -G kleister kleister

ENTRYPOINT ["/usr/bin/kleister-cli"]
CMD ["help"]

COPY --from=builder /src/bin/kleister-cli /usr/bin/kleister-cli
WORKDIR /var/lib/kleister
USER kleister
