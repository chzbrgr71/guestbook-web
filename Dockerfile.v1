FROM golang:1.8.3

MAINTAINER Brian Redmond <brianisrunning@gmail.com>

ARG VCS_REF

WORKDIR /go/src/app
COPY . .

ENV GIT_SHA $VCS_REF

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"] # ["app"]

EXPOSE 8080