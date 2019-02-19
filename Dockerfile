FROM golang:1.11.5-alpine3.8 as builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /root
COPY --from=builder /go/src/app/ .
CMD ["./main"]
