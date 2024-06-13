# Start by specifying the base image
FROM golang:1.16-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and go sum files
COPY go.mod .
COPY go.sum .

# List files in the current directory (for debugging)
RUN ls -l

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# List files in the current directory (for debugging)
RUN ls -l

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .
COPY --from=builder /app/packSizeConfig.json .

# List files in the current directory (for debugging)
RUN ls -l

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
