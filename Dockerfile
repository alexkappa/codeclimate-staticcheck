FROM golang:1.8

MAINTAINER Alex Kalyvitis <alex.kalyvitis@yieldr.com>

WORKDIR /code
VOLUME /code

COPY . /go/src/app
RUN set -x \
	&& cd /go/src/app \
	&& go build -o codeclimate-staticcheck \
	&& cp codeclimate-staticcheck /usr/local/bin

RUN useradd -u 9000 app
USER app

CMD ["/usr/local/bin/codeclimate-staticcheck"]
