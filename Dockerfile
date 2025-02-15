FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ctf-operator main.go

FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /app/ctf-operator .
CMD ["./ctf-operator"]