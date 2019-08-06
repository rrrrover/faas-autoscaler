FROM golang:1.12-alpine AS builder

WORKDIR /go/src/github.com/rrrrover/faas-autoscaler

COPY ./    ./

RUN go build -o /usr/bin/autoscaler .

FROM alpine:3.10

RUN addgroup -S app && adduser -S -g app app
RUN mkdir -p /home/app

WORKDIR /home/app

COPY --from=builder /usr/bin/autoscaler /home/app/

RUN chown -R app /home/app
USER app

ENTRYPOINT ["/home/app/autoscaler"]
