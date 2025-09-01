# Step 1: Build the Go binary
FROM golang:1.24 AS builder

# Set the working directory inside the container for the Go app
WORKDIR /app

# Copy the Go module files for caching dependencies
COPY go.mod ./

# Download Go dependencies
RUN go mod tidy

# Copy the entire Go source code into the container
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .

# Step 2: Prepare environment for the CLI and the Go app
FROM  alpine:latest

COPY --from=builder /app/myapp /root/myapp

WORKDIR /root/

# Command to run the Go binary
CMD ["./myapp"]