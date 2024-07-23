# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install git (necessary for fetching dependencies)
RUN apk update && apk add --no-cache git

# Copy the Go application source code
COPY . .

# Fetch dependencies and build the Go application
RUN  go mod tidy && go build -o main2 .

# Stage 2: Create a lightweight deployment image
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main2 .

# Make port 8080 available to the world outside this container
EXPOSE 8080

# Run the binary program we built
CMD ["./main2"]