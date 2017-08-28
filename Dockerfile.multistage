# build stage
FROM golang:latest AS build-env
MAINTAINER Brian Redmond <brianisrunning@gmail.com>
WORKDIR /app
ADD ./*.go /app/
RUN cd /app && go get github.com/denisenkom/go-mssqldb && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o guestbook-web

# final stage
FROM golang:alpine
MAINTAINER Brian Redmond <brianisrunning@gmail.com>
ARG VCS_REF
ENV GIT_SHA $VCS_REF
WORKDIR /app
COPY --from=build-env /app/guestbook-web /app/
ENTRYPOINT ./guestbook-web
EXPOSE 8080