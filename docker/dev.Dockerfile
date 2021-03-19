FROM golang:1.16-alpine
WORKDIR /go/src/github.com/logpost/jobs-optimization-service

COPY . .
RUN go get ./...
RUN go get -u github.com/cosmtrek/air

CMD air -c air.toml