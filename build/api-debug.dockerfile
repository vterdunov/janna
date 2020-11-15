FROM golang:1.14.3

ENV CGO_ENABLED 0

RUN go get github.com/derekparker/delve/cmd/dlv

WORKDIR /gomod
COPY go.mod go.sum ./
RUN go mod download

WORKDIR /go/src/github.com/vterdunov/janna

CMD [ "dlv", "debug", "./cmd/api/main.go", "--listen=:2345", "--headless=true", "--api-version=2", "--log" ]
