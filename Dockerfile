# Build stage for Go backend
FROM golang:1.21.1-alpine AS builder
# Install build dependencies
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o krtk .

# Build stage for frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY frontend/ ./
RUN npm install
RUN npm run build

# Final stage
FROM alpine:latest
WORKDIR /root/

# Install required packages
RUN apk add --no-cache sqlite-libs

# Copy backend
COPY --from=builder /app/krtk .

# Copy frontend build
COPY --from=frontend-builder /app/dist ./dist

# Copy database file
VOLUME ["/root/krtk.db"]

# Expose port
EXPOSE 8080

# Start command
CMD ["./krtk"]
