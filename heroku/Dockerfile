#######################
# BACKEND BUILD STAGE #
#######################
FROM golang:1.13-buster as golang-builder

# Get latest from GitHub
RUN mkdir -p /go/src/github.com/callummance/ && cd /go/src/github.com/callummance/ && \
    git clone https://github.com/callummance/azunyan.git

#Get golang deps and build application
RUN go get /go/src/github.com/callummance/azunyan
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build /go/src/github.com/callummance/azunyan

#######################
# FRONTEND BUILD STAGE #
#######################
RUN apt update && apt -y install nodejs npm && \
    npm config set unsafe-perm true && npm install -g yarn grunt-cli
WORKDIR /go/src/github.com/callummance/azunyan/static/frontend
RUN yarn

#############
# RUN STAGE #
#############
FROM mvertes/alpine-mongo
RUN apk add --no-cache bash

WORKDIR /root
COPY --from=golang-builder /go/azunyan .
COPY --from=golang-builder /go/src/github.com/callummance/azunyan/static/frontend ./static/frontend
COPY --from=golang-builder /go/src/github.com/callummance/azunyan/heroku/start.sh .
COPY --from=golang-builder /go/src/github.com/callummance/azunyan/azunyan.conf .
COPY --from=golang-builder /go/src/github.com/callummance/azunyan/ssh_pass.conf .
COPY --from=golang-builder /go/src/github.com/callummance/azunyan/container-scripts ./container-scripts

RUN chmod +x ./start.sh
ENTRYPOINT ./start.sh