FROM golang:1.21 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /edge-module-auto-update


FROM gcr.io/distroless/base-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /edge-module-auto-update /edge-module-auto-update
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/edge-module-auto-update"]