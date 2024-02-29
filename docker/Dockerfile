FROM golang:1.21 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /edge-module-auto-update


FROM alpine AS release-stage
WORKDIR /
COPY --from=build-stage /edge-module-auto-update /edge-module-auto-update
COPY docker/docker-compose_internal.yaml docker/docker-compose.yaml
COPY docker/config/edge_module.config docker/config/edge_module.config
COPY .env .env
EXPOSE 8080
RUN apk add docker docker-compose
USER root:root
ENTRYPOINT ["/edge-module-auto-update"]