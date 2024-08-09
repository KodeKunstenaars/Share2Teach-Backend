# Use Go 1.22.6 with Alpine Linux as the base image
FROM golang:1.22.6-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Expose the port on which the app will run
EXPOSE 8080

# Run the Go application
CMD ["go", "run", "./cmd/Share2Teach"]
