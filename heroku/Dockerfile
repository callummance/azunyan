FROM golang:1.13-buster

# Install mongodb
RUN apt update && apt -y install gnupg2 && \
    wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | apt-key add -

RUN echo "deb http://repo.mongodb.org/apt/debian buster/mongodb-org/4.2 main" | tee /etc/apt/sources.list.d/mongodb-org.list

RUN apt update && apt -y install mongodb-org

RUN apt -y install nodejs npm && npm config set unsafe-perm true && npm install -g yarn grunt-cli 

#Get the latest version from git
RUN go get -u github.com/dennisTsiang/azunyan

#Include files from host
ADD . /go/src/github.com/dennisTsiang/azunyan

#Get golang deps
RUN go get /go/src/github.com/dennisTsiang/azunyan
RUN go install -i github.com/dennisTsiang/azunyan

#Change pwd so that static files work fine
WORKDIR /go/src/github.com/dennisTsiang/azunyan/static/frontend

#Get frontend deps
RUN yarn

WORKDIR /go/src/github.com/dennisTsiang/azunyan/
RUN chmod +x heroku/start.sh

ENTRYPOINT ./heroku/start.sh