# Start from the official Go image
FROM golang:1.23.4-alpine

# Set environment variables
ENV GO111MODULE=on 
ENV PORT=3000

# Install necessary packages
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o app ./cmd

# Expose the API port
EXPOSE 3000

# Run the application
CMD ["./app"]
