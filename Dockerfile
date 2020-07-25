FROM golang:1.13-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache --update bash git gcc g++ && \
    go get -u -v github.com/kardianos/govendor

WORKDIR /go/src/github.com/guillaumejacquart/go-http-scheduler
COPY . .

RUN GOOS=linux go build -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/guillaumejacquart/go-http-scheduler/app .
COPY --from=0 /go/src/github.com/guillaumejacquart/go-http-scheduler/config.yml .
COPY --from=0 /go/src/github.com/guillaumejacquart/go-http-scheduler/public ./public

EXPOSE 8080

CMD ["./app"]  