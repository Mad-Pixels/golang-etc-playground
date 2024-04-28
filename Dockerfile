ARG ALPINE_VERSION=3.18
ARG GO_VERSION=1.21
ARG NAME=bootstrap

ARG GOCACHE=/root/.cache/go-build
ARG ASM_FLAGS="-trimpath"
ARG GC_FLAGS="-trimpath"
ARG LD_FLAGS="-w -s"
ARG GOOS=linux

# amd64
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-amd64
ARG ASM_FLAGS
ARG GC_FLAGS
ARG LD_FLAGS
ARG GOCACHE
ARG NAME

COPY apps /go/src/apps
WORKDIR /go/src/apps

ENV GOARCH=amd64
RUN go mod vendor
RUN --mount=type=cache,target="${GOCACHE}" \
      go build \
        -mod=vendor \
        -asmflags="${ASM_FLAGS}" \
        -ldflags="${LD_FLAGS}"   \
        -gcflags="${GC_FLAGS}"   \
        -o /out/${NAME}          \
      /go/src/apps/${NAME}/cmd

FROM alpine:${ALPINE_VERSION} AS amd64
ARG NAME
ENV APP=${NAME}
COPY --from=builder-amd64 /out/${APP} /${APP}
ENTRYPOINT ["sh", "-c", "/$APP"]

# arm64
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-arm64
ARG ASM_FLAGS
ARG GC_FLAGS
ARG LD_FLAGS
ARG GOCACHE
ARG NAME

COPY apps /go/src/apps
WORKDIR /go/src/apps

ENV GOARCH=arm64
RUN go mod vendor
RUN --mount=type=cache,target="${GOCACHE}" \
      go build \
        -mod=vendor \
        -asmflags="${ASM_FLAGS}" \
        -ldflags="${LD_FLAGS}"   \
        -gcflags="${GC_FLAGS}"   \
        -o /out/${NAME}          \
      /go/src/apps/${NAME}/cmd

FROM alpine:${ALPINE_VERSION} AS arm64
ARG NAME
ENV APP=${NAME}
COPY --from=builder-arm64 /out/${APP} /${APP}
ENTRYPOINT ["sh", "-c", "/$APP"]
