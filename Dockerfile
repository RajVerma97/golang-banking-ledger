FROM golang:1.23.4-alpine

ENV GO111MODULE=on 
ENV PORT=3000

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd

EXPOSE 3000

CMD ["./app"]
