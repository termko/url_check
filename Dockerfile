FROM golang:1.16.3

WORKDIR  /app

COPY . .

RUN go build 

CMD ["./ozon_service"]
