FROM golang:1.8

MAINTAINER Alex Kalyvitis <alex.kalyvitis@yieldr.com>

WORKDIR /code
VOLUME /code

COPY bin/codeclimate-staticcheck-linux-amd64 /usr/local/bin/codeclimate-staticcheck

# RUN set -x \
# 	&& cd /go/src/app \
# 	&& go build -o codeclimate-staticcheck \
# 	&& cp codeclimate-staticcheck /usr/local/bin

RUN set -x \
	&& useradd -u 9000 app \
	&& chown -R app /go/src

USER app

CMD ["/usr/local/bin/codeclimate-staticcheck"]
