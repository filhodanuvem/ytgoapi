FROM golang:1.21-alpine3.18 as builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/api
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o api .

FROM scratch
WORKDIR /
COPY --from=builder /app/cmd/api ./
EXPOSE 3000
ENTRYPOINT ["./api"]