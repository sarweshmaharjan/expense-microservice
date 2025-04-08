# Use the official Golang image as the base image
FROM golang:1.22.2

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . ./view

# Set the working directory to the cmd folder
WORKDIR /app/view/cmd

RUN go mod tidy \
    && go mod verify

# Build the Go application
RUN go build -o /main

# Command to run the executable
CMD ["/main"]