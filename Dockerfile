FROM --platform=$BUILDPLATFORM golang:1.20.3-alpine3.17 AS builder

WORKDIR /build

# Prefetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/elab-backend cmd/main.go

FROM alpine:3.17

RUN apk add --no-cache ca-certificates
COPY ./configs/config.toml /app/configs/config.toml
COPY --from=builder /app/elab-backend /app/elab-backend

WORKDIR /app

ENTRYPOINT [ "./elab-backend" ]
