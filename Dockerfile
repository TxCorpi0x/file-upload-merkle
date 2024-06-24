FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o ./fxmerkle ./main.go

FROM alpine
EXPOSE 8080
COPY --from=builder /app/fxmerkle /fxmerkle
CMD ["/fxmerkle", "server"]
