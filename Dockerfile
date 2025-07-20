# STAGE 1: build
FROM golang:1.24.5-alpine AS build
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -trimpath -ldflags "-s -w" -o /app/bin/simple-marketplace ./cmd/api

# STAGE 2: runtime
FROM alpine:3.22.1
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=build /app/bin/simple-marketplace /app/simple-marketplace
EXPOSE 8080
USER app
ENTRYPOINT ["/app/simple-marketplace"]