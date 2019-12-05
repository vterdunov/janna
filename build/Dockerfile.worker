FROM golang:1.13.4-alpine3.10

ENV CGO_ENABLED 0

RUN apk add --no-cache git
RUN go get github.com/derekparker/delve/cmd/dlv

FROM alpine:3.10.1

EXPOSE 2345

WORKDIR /janna

COPY --from=0 /go/bin/dlv /
COPY worker worker

CMD ["./worker"]
# CMD [ "/dlv", "exec", "./worker", "--listen=:2345", "--headless=true", "--api-version=2", "--log" ]
