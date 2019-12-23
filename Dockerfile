#######################
# BACKEND BUILD STAGE #
#######################
FROM golang:1.13.5-buster as golang-builder

#Include files from host
ADD . /go/src/github.com/callummance/azunyan

#Get golang deps
RUN go get /go/src/github.com/callummance/azunyan
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o azunyan /go/src/github.com/callummance/azunyan



########################
# FRONTEND BUILD STAGE #
########################
FROM node:lts-alpine as nodejs-builder

ADD ./static/frontend /home/node/app
WORKDIR /home/node/app
RUN yarn



#############
# RUN STAGE #
#############
FROM alpine:latest
RUN apk add --no-cache bash

WORKDIR /root
COPY --from=golang-builder /go/azunyan .
COPY --from=nodejs-builder /home/node/app ./static/frontend

ENTRYPOINT /root/azunyan -c /run/secrets/azunyan_conf

EXPOSE 8080
