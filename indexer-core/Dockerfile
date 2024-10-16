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

# Install migrate tool (adjust version as needed)
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz && \
    tar xzf migrate.linux-amd64.tar.gz && \
    mv migrate /usr/local/bin/ && \
    rm migrate.linux-amd64.tar.gz

# Expose the port the app runs on
EXPOSE 8080

# Command to run Air with migration
CMD ["sh", "-c", "echo 'Running migrations...' && migrate -path migrations -database \"$PG_URL\" up && echo 'Migrations completed!' && air"]
