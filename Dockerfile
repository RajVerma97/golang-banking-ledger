FROM golang:1.23.4-alpine AS builder

ENV GO111MODULE=on 

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 3000

CMD ["./app"]