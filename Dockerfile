# Start from the official Go image
FROM golang:1.17-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o myapp

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory to /app in the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
