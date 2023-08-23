# Build stage
FROM golang:1.20-alpine as builder
RUN mkdir /build
WORKDIR /build
COPY . .
ENV GOOS=linux CGO_ENABLED=0
RUN set -ex && \
    apk add --no-progress --no-cache \
    gcc \
    musl-dev

RUN go build -o server ./main.go

# Run stage
FROM alpine:3.16
RUN apk --no-cache add ca-certificates ffmpeg
WORKDIR /app
COPY --from=builder /build/server .
COPY .env .
COPY start.sh .
COPY wait-for-it.sh .
COPY migration ./migration
RUN chmod +x start.sh wait-for-it.sh

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/server"]