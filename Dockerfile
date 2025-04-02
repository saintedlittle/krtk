FROM golang:1.21.1-alpine AS builder
# Install build dependencies
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o krtk .

FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache sqlite-libs
COPY --from=builder /app/krtk .
COPY templates/ ./templates/
RUN chmod +x krtk
CMD ["./krtk"]