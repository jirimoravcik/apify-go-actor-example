FROM golang:1.17.1-alpine

WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o /example-actor
CMD ["/example-actor"]
