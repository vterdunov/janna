ARG GOLANG_VERSION=1.15.5
# ---------------
FROM golang:${GOLANG_VERSION} AS env

RUN mkdir /user && \
    echo 'janna:x:65534:65534:janna:/:' > /user/passwd && \
    echo 'janna:x:65534:' > /user/group

WORKDIR /gomod
COPY go.mod go.sum ./
RUN go mod download

# ---------------
FROM env AS builder
WORKDIR /build
COPY . .

RUN make compile

# ---------------
FROM builder AS tester

RUN make test

# ---------------
FROM builder AS licence-scanner

ARG GOLICENSE_VERSION=0.1.1
RUN wget https://github.com/mitchellh/golicense/releases/download/v${GOLICENSE_VERSION}/golicense_${GOLICENSE_VERSION}_linux_x86_64.tar.gz && \
    tar xvf golicense_${GOLICENSE_VERSION}_linux_x86_64.tar.gz
ARG GITHUB_TOKEN
RUN ./golicense -verbose ./build/ci/golicense.hcl ./janna

# ---------------
FROM alpine:3.9.3 AS final

ENV PORT=8080

WORKDIR /janna

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER janna:janna

CMD ["./janna"]

COPY --from=builder /build/janna janna
