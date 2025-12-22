# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application as a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -o main ./cmd/app

# Create a lightweight image for the final binary
FROM gcr.io/distroless/static

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose application port (change if needed)
ENV PORT=8080
EXPOSE ${PORT}

# Run the application
CMD ["./main"]