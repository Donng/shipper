FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app/consignment-service

COPY . .

ENV HTTP_PROXY "此处需要配合代理地址"
ENV HTTPS_PROXY "此处需要配置代理地址"

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consignment-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/consignment-service/consignment-service .

CMD ["./consignment-service"]
