FROM alpine:3.21.3
WORKDIR /app

RUN apk add --no-cache tzdata

COPY domain-watch .

ARG USERNAME=domain-watch
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"
USER $UID

CMD ["./domain-watch"]
