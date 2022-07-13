FROM golang:1.18.1-alpine3.15 AS builder
RUN apk update
RUN apk add git
RUN mkdir /build
ADD go.mod /build/
ADD go.sum /build/
ADD ./app /build/app
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app-binary ./app/

FROM alpine:3.15.4
COPY --from=builder /build/app-binary /app
WORKDIR /
RUN chmod +x ./app
CMD ["./app"]
# Expose port for HTTP server
EXPOSE 8008
