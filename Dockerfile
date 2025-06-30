FROM ubuntu:latest
LABEL authors="aleksandrkozlov"

ENTRYPOINT ["top", "-b"]