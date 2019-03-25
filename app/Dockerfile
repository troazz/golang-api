FROM golang:alpine

ADD ./src /go/src/app

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go get -u "github.com/gin-gonic/gin"
RUN go get -u "github.com/go-sql-driver/mysql"
RUN go get -u "github.com/streadway/amqp"
RUN go get -u "github.com/olivere/elastic"
RUN go get "github.com/patrickmn/go-cache"

WORKDIR /go/src/app
ENV PORT=3001

CMD ["go", "run", "main.go"]