FROM golang:1.17.2-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN go build -o ttsaas

FROM alpine:3.12

RUN apk update \
    && apk add bash \
    && apk add ffmpeg \
    && apk add mplayer \
    && rm -rf /var/chache/apk/* \
    && addgroup -S app && adduser -S app -G app

USER app

WORKDIR /app

COPY --from=builder /app/ttsaas .

ENTRYPOINT ["./ttsaas"]
