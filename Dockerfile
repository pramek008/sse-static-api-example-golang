# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY main.go .

# Initialize go mod (karena kita hanya pakai stdlib, ini formalitas agar rapi)
RUN go mod init llm-stream-go

# Build binary
# -ldflags="-s -w" membuang debug info agar ukuran file lebih kecil
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server main.go

# Stage 2: Run (Tiny Image)
FROM alpine:latest

# Security: Buat user non-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy hanya binary dari stage 1
COPY --from=builder /app/server .

# Gunakan user non-root
USER appuser

# Expose port
EXPOSE 3950

# Jalankan server
CMD ["./server"]