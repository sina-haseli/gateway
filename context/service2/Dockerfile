# Use an official Go image as a base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app


# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY ../../cmd/service2_mock ./cmd/service2_mock
COPY ../../pkg ./pkg

# Build the Go app
RUN go build -o service2_mock cmd/service2_mock/main.go

# Expose port 8080
EXPOSE 50051

# Command to run the executable
CMD ["./service2_mock"]
