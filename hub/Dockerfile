# Use the official Go image as the base image
FROM golang:1.22.2 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . ./hub

WORKDIR /app/hub/cmd/

# Build the Go application
RUN go build -o /main

# Command to run the executable
CMD ["/main"]
