FROM golang:1-alpine as builder

WORKDIR /app

COPY go.mod .
RUN go mod download

ARG GOOS=linux
ARG GOARCH=amd64
COPY . .
RUN go build -ldflags="-w -s"

FROM alpine

WORKDIR /app

ENV PATH="/app:$PATH"

COPY --from=builder /app/domain-expiration-notifier .
ENTRYPOINT ["domain-expiration-notifier"]
