# Stage 1: Build the Go binary
# This stage compiles your Go code into a binary
FROM golang:1.21-alpine AS builder

# Set working directory inside container
WORKDIR /app

# Copy go mod files
COPY go.mod go.mod

# Copy all source code
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build the binary
# CGO_ENABLED=0 means it's statically compiled (works anywhere)
# -o cost-detector means output name
RUN CGO_ENABLED=0 GOOS=linux go build -o cost-detector ./cmd/cost-detector

echo "Build complete! Binary created."

# Stage 2: Create the runtime image
# Use alpine (tiny Linux image) for the final container
FROM alpine:3.18

# Install ca-certificates so we can talk to HTTPS endpoints
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy ONLY the binary from the builder stage
# This keeps the image small (no Go compiler, no source code)
COPY --from=builder /app/cost-detector .

# Expose a port (optional, for health checks)
EXPOSE 8080

# Run the binary when container starts
CMD ["./cost-detector"]