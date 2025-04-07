# Stage 1: Build the application
FROM golang:1.23.6-alpine AS builder

WORKDIR /usr/local/src

# Install necessary packages
RUN apk --no-cache add bash git make gcc musl-dev

# Copy go.mod and go.sum files
COPY ["go.mod", "go.sum", "./"]

# Download dependencies
RUN go mod download

# Copy the source code
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

# Build the application
RUN go build -o /usr/local/bin/app cmd/main/main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest AS runner

# Copy the built application from the builder stage
COPY --from=builder /usr/local/bin/app /app

# Copy the configuration file
COPY config.yml ./config.yml

# Set the entrypoint command
CMD ["/app" ]