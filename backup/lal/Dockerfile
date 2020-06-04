# Builder stage
FROM golang:alpine AS builder

ADD . /go/src/app
WORKDIR /go/src/app

CMD ["go", "run", "main.go"]