FROM golang:1.21.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o krtk ./cmd/krtk

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/krtk .
RUN chmod +x krtk
CMD ["./krtk"]