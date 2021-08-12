FROM alpine:latest

ADD out/bin/wow-server /usr/local/bin/wow-server

RUN chmod u+x /usr/local/bin/wow-server

ENTRYPOINT ["wow-server"]
