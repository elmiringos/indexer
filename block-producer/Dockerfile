# Use the official Golang image
FROM golang:1.22 AS dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Command to run Air with migration
CMD ["sh", "-c", "air"]
