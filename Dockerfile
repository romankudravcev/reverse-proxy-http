# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies first (if any)
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o proxy

# Final stage
FROM golang:1.23 AS RUNTIME

# Copy the binary from the builder stage
COPY --from=builder /app/proxy /proxy

# Expose the desired port
EXPOSE 8734

# Set the entrypoint
CMD ["/proxy"]
