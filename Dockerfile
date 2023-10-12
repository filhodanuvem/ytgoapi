FROM golang:1.21-alpine3.18

WORKDIR /app

COPY . .

RUN go install github.com/cosmtrek/air@latest

RUN go get -d -v ./...

EXPOSE 3000

CMD ["air"]