# Use the official Golang image
FROM golang:1.23 AS dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

CMD ["go", "run", "cmd/api/main.go"]
