# syntax=docker/dockerfile:1

# Build phase
FROM golang:1.21 AS builder
LABEL authors="gerhard"

# Set the working directory inside the build container
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app (CGO disabled for a fully static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o share2teach-api ./cmd/api/

# Final phase (using Alpine for minimal image size)
FROM alpine:latest

# Install certificates (optional, if your app requires HTTPS requests)
RUN apk --no-cache add ca-certificates

# Set the working directory inside the Alpine container
WORKDIR /root/

# Copy the binary from the build container
COPY --from=builder /app/share2teach-api .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./share2teach-api"]