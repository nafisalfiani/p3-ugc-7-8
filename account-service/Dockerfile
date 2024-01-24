# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go app
RUN go build -o ./build/app ./main.go

# Create blank env
RUN touch .env

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./build/app"]
