FROM golang:1.21.0-alpine3.17
WORKDIR /app
COPY . /app/
RUN go build -o main main.go

EXPOSE 8080


