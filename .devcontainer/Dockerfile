FROM mcr.microsoft.com/devcontainers/base:ubuntu

ARG GO_VERSION="1.26.0"
ARG GOLANGCI_LINT_VERSION="v2.10.1"
ARG GOFUMPT_VERSION="v0.9.2"
ARG GOPLS_VERSION="v0.21.1"
ARG DLV_VERSION="v1.26.0"
ARG NATSCLI_VERSION="v0.3.1"

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends curl ca-certificates build-essential \
    && ARCH=$(dpkg --print-architecture) \
    && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" -o /tmp/go.tar.gz \
    && mkdir -p /usr/local/go \
    && tar -xzf /tmp/go.tar.gz -C /usr/local --strip-components=1 \
    && rm /tmp/go.tar.gz \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

ENV GOPATH=/home/vscode/go
ENV PATH=$PATH:/usr/local/go/bin:${GOPATH}/bin

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}

USER vscode

RUN CGO_ENABLED=0 go install mvdan.cc/gofumpt@${GOFUMPT_VERSION} \
    && CGO_ENABLED=0 go install golang.org/x/tools/gopls@${GOPLS_VERSION} \
    && CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@${DLV_VERSION} \
    && CGO_ENABLED=0 go install github.com/nats-io/natscli/nats@${NATSCLI_VERSION}