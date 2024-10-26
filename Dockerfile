FROM golang:1.22.5-alpine as builder

WORKDIR /go/src/github.com/corka149/rental

COPY . .

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rental ./cmd/rental/rental.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/github.com/corka149/rental/rental .

RUN chmod +x rental

CMD ["./rental"]
