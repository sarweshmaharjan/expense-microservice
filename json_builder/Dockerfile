# Use the official Golang image as the base image
FROM golang:1.22.2

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./

RUN go mod download

# Copy the source code into the container
COPY . ./json_builder

# Set the working directory to the main.go location
WORKDIR /app/json_builder/cmd

# Build the Go application
RUN go build -o /main

# Command to run the executable
CMD ["/main"]
