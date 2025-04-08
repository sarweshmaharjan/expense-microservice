# Use the official Golang image as the base image
FROM golang:1.22.2

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . ./report_generator

# Set the working directory to the main.go directory
WORKDIR /app/report_generator/cmd

# Build the Go application
RUN go build -o /main

# Command to run the executable
CMD ["/main"]
