FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go sidecar.go ./

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o sidecar sidecar.go

# Create a minimal production image
FROM alpine:3.18

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/sidecar /app/sidecar

# Run the main application by default
ENTRYPOINT ["/app/main"]