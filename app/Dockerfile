FROM golang:1.13 AS builder
WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 go build main.go

FROM alpine:3.10
COPY --from=builder /go/src/main .
RUN chmod a+x ./main
CMD ["./main"]