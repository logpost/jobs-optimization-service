FROM golang:1.16-alpine

ARG GIT_ACCESS_TOKEN_CURL_CONFIG

COPY . /go/src/github.com/logpost/jobs-optimization-service
WORKDIR /go/src/github.com/logpost/jobs-optimization-service

RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/logpost/logpost-environment/master/environment/jobs-optimization-service/config.toml
RUN mkdir conf && mv -f production.yml staging.yml development.yml conf && mkdir builder
RUN go get ./...
RUN go build -o ./builder/jobs-optimization-svc

EXPOSE 8083
ENV KIND=stagging 

CMD ["./builder/jobs-optimization-svc"]

