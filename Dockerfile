FROM golang
MAINTAINER Brian Redmond <brianisrunning@gmail.com>
# full golang: 727MB
# multi-stage: 256MB

ARG VCS_REF
ENV GIT_SHA $VCS_REF

WORKDIR /app
ADD ./*.go /app/

RUN cd /app && go get github.com/denisenkom/go-mssqldb && go build -o guestbook-web

ENTRYPOINT /app/guestbook-web
EXPOSE 8080