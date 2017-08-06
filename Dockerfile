FROM alpine:3.4
MAINTAINER Alex Kalyvitis <alex.kalyvitis@yieldr.com>
COPY bin/codeclimate-staticcheck-linux-amd64 /usr/local/bin/codeclimate-staticcheck
WORKDIR /code
VOLUME /code
RUN adduser -u 9000 -D app
USER app
CMD ["/usr/local/bin/codeclimate-staticcheck"]
