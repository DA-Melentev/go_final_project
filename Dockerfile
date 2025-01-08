FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /task-app cmd/server/main.go
CMD ["/task-app"]
