FROM alpine:3.4
MAINTAINER Alex Kalyvitis <alex.kalyvitis@yieldr.com>
WORKDIR /usr/src/app
COPY codeclimate-staticcheck.go /usr/src/app/codeclimate-staticcheck.go
ADD engine.json /
RUN apk --update add go git && \
  export GOPATH=/tmp/go GOBIN=/usr/local/bin && \
  go get -d . && \
  go install codeclimate-staticcheck.go && \
  apk del go git && \
  rm -rf "$GOPATH" && \
  rm /var/cache/apk/*
WORKDIR /code
VOLUME /code
RUN adduser -u 9000 -D app
USER app
CMD ["/usr/local/bin/codeclimate-staticcheck"]
