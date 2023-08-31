#Build stage
FROM golang:1.21.0-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Run stage
#below image is the linux alpine image(lightweight img)
FROM alpine:3.17
WORKDIR /app
# . refers to working dir in the base image which is /app here.
COPY --from=builder /app/main .
COPY app.env .       
EXPOSE 8080
CMD ["/app/main"]
            
