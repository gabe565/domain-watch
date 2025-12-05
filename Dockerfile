FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.9.0 AS xx

FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS builder
WORKDIR /app

COPY --from=xx / /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETPLATFORM
RUN --mount=type=cache,target=/root/.cache \
    CGO_ENABLED=0 xx-go build -ldflags='-w -s' -trimpath


FROM alpine:3.21.3
WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=builder /app/domain-watch /usr/local/bin/

ARG USERNAME=domain-watch
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"
USER $UID

ENTRYPOINT ["domain-watch"]
