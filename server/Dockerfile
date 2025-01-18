# Use a lightweight Go image
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy project files
# COPY go.mod go.sum ./
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Start a new minimal image for running the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy other necessary folders (e.g., routes, uploads, tmp)
COPY --from=builder /app/smtp /app/smtp
# COPY --from=builder /app/lib /app/lib

COPY --from=builder /app/.env .

# Expose SMTP ports
EXPOSE 25 587 2525

# Command to run the server
CMD ["./main"]
