FROM golang:1.16-buster as builder
WORKDIR /go/src/github.com/logpost/jobs-optimization-service

ARG GIT_ACCESS_TOKEN_CURL_CONFIG

COPY go.* ./
RUN go mod download
COPY . .

RUN go build -mod=readonly -v -o ./jobs-optimization-svc 

FROM debian:buster-slim
WORKDIR /go/src/github.com/logpost/jobs-optimization-service/jobs-optimization-svc
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/logpost/jobs-optimization-service/jobs-optimization-svc .
COPY --from=builder /go/src/github.com/logpost/jobs-optimization-service/conf conf/

ENV GO111MODULE=on
ENV PORT=${PORT}
ENV KIND=production

EXPOSE 8083

CMD ["./jobs-optimization-svc"]
