FROM golang:1.16
WORKDIR /go/src/github.com/logpost/jobs-optimization-service

ARG GIT_ACCESS_TOKEN_CURL_CONFIG

COPY . /go/src/github.com/logpost/jobs-optimization-service

RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/logpost/logpost-environment/master/environment/jobs-optimization-service/config.toml
RUN mkdir conf && mv -f config.toml conf
RUN go get ./...
RUN go build -mod=readonly -v -o ./jobs-optimization-svc
 

EXPOSE 8083
ENV KIND=stagging 

CMD ["./jobs-optimization-svc"]
