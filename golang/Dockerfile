# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code and config
COPY . .

# Build both applications - fixing the path to match your actual directory structure
RUN CGO_ENABLED=0 GOOS=linux go build -o main-server cmd/server/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o grpc-server cmd/grcp-server/main.go 
    # Note: Fixed path to match your directory (grcp-server)

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install bash and necessary certificates
RUN apk add --no-cache bash ca-certificates

# Copy the binaries and config from builder
COPY --from=builder /app/main-server .
COPY --from=builder /app/grpc-server .
COPY --from=builder /app/config.yaml .

# Copy migrations directory with correct path
COPY --from=builder /app/ecomm-grpc/db/migrations /app/db/migrations/

# Create and configure startup script
RUN printf '#!/bin/bash\n\
# Wait for database\n\
echo "Waiting for database..."\n\
sleep 10\n\
\n\
# Start gRPC server in background\n\
echo "Starting gRPC server..."\n\
./grpc-server & \n\
GRPC_PID=$!\n\
\n\
# Wait for gRPC server to start\n\
echo "Waiting for gRPC server to initialize..."\n\
sleep 5\n\
\n\
# Start HTTP server\n\
echo "Starting HTTP server..."\n\
./main-server & \n\
HTTP_PID=$!\n\
\n\
# Wait for either process to exit\n\
wait -n $GRPC_PID $HTTP_PID\n\
\n\
# If one process exits, kill the other\n\
kill $(jobs -p) 2>/dev/null\n' > /app/start.sh && \
    chmod +x /app/start.sh

# Create a non-root user
RUN adduser -D appuser && \
    chown -R appuser:appuser /app && \
    chmod -R 755 /app

USER appuser

EXPOSE 8888 50051

CMD ["/app/start.sh"]