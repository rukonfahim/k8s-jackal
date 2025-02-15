FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o ctf-operator main.go
FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /app/ctf-operator .
CMD ["./ctf-operator"]