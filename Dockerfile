# syntax=docker/dockerfile:1

# Base stage for shared configurations
FROM golang:1.23.3-alpine AS base
WORKDIR /app
RUN apk add --no-cache gcc musl-dev make

# Development stage
FROM base AS development
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build
# Install CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon@latest

# Create build directory and set permissions
RUN mkdir -p /app/build && \
    chown -R 1000:1000 /app && \
    chmod -R 777 /app

EXPOSE 8080


# Production stage
FROM alpine:latest AS production
RUN apk --no-cache add make

WORKDIR /app
COPY --from=base /app/build/url-shortener .
COPY --from=base /app/Makefile .

EXPOSE 8080

ENTRYPOINT ["./url-shortener"]