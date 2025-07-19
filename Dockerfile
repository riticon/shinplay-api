# Use the official Golang image as a base image
FROM golang:1.24.2

# Set the working directory inside the container
WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Create tmp directory for Air
RUN mkdir -p tmp

# Expose the application port
EXPOSE 8080

# Use Air for hot reload in development
CMD ["air"]
