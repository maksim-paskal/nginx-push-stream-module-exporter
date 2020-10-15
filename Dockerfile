FROM golang:1.14 as build

COPY ./cmd /usr/src/nginx-push-stream-module-exporter/cmd
COPY go.* /usr/src/nginx-push-stream-module-exporter/
COPY .git /usr/src/nginx-push-stream-module-exporter/

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOFLAGS="-trimpath"

RUN cd /usr/src/nginx-push-stream-module-exporter \
  && go mod download \
  && go mod verify \
  && go build -v -o nginx-push-stream-module-exporter -ldflags "-X main.buildTime=$(date +"%Y%m%d%H%M%S") -X main.gitVersion=`git describe --exact-match --tags $(git log -n1 --pretty='%h')`" ./cmd

FROM alpine:latest

COPY --from=build /usr/src/nginx-push-stream-module-exporter/nginx-push-stream-module-exporter /app/nginx-push-stream-module-exporter

WORKDIR /app

RUN addgroup -g 101 -S app \
&& adduser -u 101 -D -S -G app app

USER 101

CMD /app/nginx-push-stream-module-exporter