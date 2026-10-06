# Multi-stage build (docs/09-devops.md §1): binary statis ±20 MB, non-root.

FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/fablerace-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app
USER app
COPY --from=build /out/fablerace-api /usr/local/bin/fablerace-api
EXPOSE 8080
ENTRYPOINT ["fablerace-api"]
