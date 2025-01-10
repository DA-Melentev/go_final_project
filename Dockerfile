FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /task-app cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /task-app .

CMD ["./task-app"]
