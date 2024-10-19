FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o proxy

EXPOSE 8734

CMD ["./proxy"]

