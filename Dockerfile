#Build stage
FROM golang:1.21.0-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add wget
RUN wget -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz
RUN tar xvfz migrate.linux-amd64.tar.gz


#Run stage
#below image is the linux alpine image(lightweight img)
FROM alpine:3.17
WORKDIR /app
# . refers to working dir in the base image which is /app here.
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY ./db/migration/ /app/migration/
COPY app.env .
COPY ./start.sh /app/       
EXPOSE 8080
CMD ["/app/main"]
#below cmd is to frst execute the start.sh(to migrate db up) and start app
ENTRYPOINT [ "/app/start.sh" ]
            
