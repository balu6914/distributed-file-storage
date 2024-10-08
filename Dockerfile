FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main .
RUN chmod +x /app/main

CMD ["/app/main"]
