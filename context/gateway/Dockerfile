# Use an official Go image as a base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app


# Copy go mod and sum files
COPY ../../go.mod ./
COPY ../../go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY ../.. .

# Build the Go app
RUN go build -o gateway cmd/server/main.go

# Expose port 8080
EXPOSE 9090

# Command to run the executable
CMD ["./gateway"]
