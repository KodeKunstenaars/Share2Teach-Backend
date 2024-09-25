# syntax=docker/dockerfile:1
FROM golang:1.21
LABEL authors="gerhard"

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend directory, including internal packages
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /share2teach-api ./cmd/api/

# Expose port
EXPOSE 8080

# Run the binary
CMD ["/share2teach-api"]