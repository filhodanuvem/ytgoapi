FROM golang:1.21-alpine3.18

WORKDIR /app

COPY . .

WORKDIR /app/cmd/api

RUN go get -d -v ./...

RUN go build -o api .

EXPOSE 3000

CMD ["./api"]
