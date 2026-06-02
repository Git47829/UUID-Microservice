FROM golang:latest AS build-stage

WORKDIR /app

COPY go.sum go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -p /uuid-service

FROM gcr.io/distroless/base-debian11 as build-release

WORKDIR /

COPY --from=build-stage /uuid-service /uuid-service

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/uuid-service"]
