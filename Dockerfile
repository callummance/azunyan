#Use the official golang image based on alpine
FROM golang:1.11-alpine

#Need git to fetch packages
RUN apk update && apk upgrade && apk add --no-cache bash git openssh

#Also install yarn to get js packages for frontend
RUN apk add --no-cache nodejs npm && \
    npm install -g yarn grunt-cli && \

#Get the latest version from git
RUN go get -u github.com/callummance/azunyan

#Include files from host
ADD . /go/src/github.com/callummance/azunyan

#Get golang deps
RUN go get /go/src/github.com/callummance/azunyan
RUN go install -i github.com/callummance/azunyan

#Change pwd so that static files work fine
WORKDIR /go/src/github.com/callummance/azunyan/static/frontend

#Get frontend deps
RUN yarn

WORKDIR /go/src/github.com/callummance/azunyan/

ENTRYPOINT /go/bin/azunyan -c /run/secrets/azunyan_conf

EXPOSE 8080
