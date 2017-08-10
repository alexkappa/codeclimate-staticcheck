FROM debian
MAINTAINER Alex Kalyvitis <alex.kalyvitis@yieldr.com>

WORKDIR /code
VOLUME /code

COPY bin/codeclimate-staticcheck-linux-amd64 /usr/local/bin/codeclimate-staticcheck

RUN useradd -u 9000 app
USER app

CMD ["/usr/local/bin/codeclimate-staticcheck"]
