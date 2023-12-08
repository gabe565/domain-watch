FROM --platform=$BUILDPLATFORM golang:1.21-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Set Golang build envs based on Docker platform string
ARG TARGETPLATFORM
RUN --mount=type=cache,target=/root/.cache \
    set -x \
    && case "$TARGETPLATFORM" in \
        'linux/amd64') export GOARCH=amd64 ;; \
        'linux/arm/v6') export GOARCH=arm GOARM=6 ;; \
        'linux/arm/v7') export GOARCH=arm GOARM=7 ;; \
        'linux/arm64') export GOARCH=arm64 ;; \
        *) echo "Unsupported target: $TARGETPLATFORM" && exit 1 ;; \
    esac \
    && go build -ldflags='-w -s' -trimpath


FROM alpine:3.19
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
