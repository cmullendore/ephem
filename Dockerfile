FROM ubuntu:latest

ADD ephem /
WORKDIR /local
EXPOSE 8443/tcp
ENTRYPOINT ["/ephem"]