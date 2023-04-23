FROM golang:1.20-alpine as builder
RUN apk update && apk add curl openssl git openssh-client build-base
ADD . /src
WORKDIR /src
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /src/dist/vault-backup

FROM alpine:latest as done
COPY --from=builder /src/dist /opt
WORKDIR /opt
ENTRYPOINT ["./vault-backup"]