FROM golang:onbuild
LABEL maintainer "maintainer@spirala.co" description "Management UI for Docker Swarm"

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o spirala .
RUN ["/app/spirala"]
