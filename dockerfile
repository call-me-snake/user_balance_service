FROM golang:1.12.0-alpine3.9 AS builder
WORKDIR /go/src/github.com/call-me-snake/user_balance_service
COPY . .
RUN go install ./...

FROM jwilder/dockerize AS production
COPY --from=builder /go/bin/cmd ./app

#docker build -t user_balance_service_img .
#docker run -it --name balance_service user_balance_service_img /bin/sh
#docker stop balance_service
#docker rm balance_service