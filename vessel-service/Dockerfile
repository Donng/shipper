FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app/vessel-service

COPY . .

ENV HTTP_PROXY "修改为自己的代理服务"
ENV HTTPS_PROXY "修改为自己的代理服务"

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vessel-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/vessel-service .

CMD ["./vessel-service"]
