# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.23-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/abcsun-qr .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S abcsun && adduser -S -G abcsun abcsun && mkdir -p /data && chown abcsun:abcsun /data

COPY --from=build /out/abcsun-qr /usr/local/bin/abcsun-qr

WORKDIR /data
USER abcsun

ENV PORT=18080
EXPOSE 18080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
	CMD wget -qO- http://127.0.0.1:18080/api/v1/readyz >/dev/null || exit 1

ENTRYPOINT ["abcsun-qr"]
