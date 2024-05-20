# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o /app/main ./api/main.go

# Stage 2: Create a minimal image with the built application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/main .

# Expose the port the application runs on (adjust if necessary)
EXPOSE 3000

# Command to run the application
CMD ["./main"]
