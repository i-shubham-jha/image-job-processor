# Use the official Golang image as a build stage
FROM golang:1.23.1

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o image-job-processor ./cmd/image-job-processor/main.go
RUN chmod +x image-job-processor

# Command to run the executable
CMD ["./image-job-processor"]
