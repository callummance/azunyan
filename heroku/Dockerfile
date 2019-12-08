FROM golang:1.13-buster

# Install mongodb
RUN apt update && apt -y install gnupg2 && \
    wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | apt-key add -

RUN echo "deb http://repo.mongodb.org/apt/debian buster/mongodb-org/4.2 main" | tee /etc/apt/sources.list.d/mongodb-org.list

RUN apt update && apt -y install mongodb-org

RUN apt -y install nodejs npm && npm config set unsafe-perm true && npm install -g yarn grunt-cli 

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
RUN chmod +x heroku/start.sh

ENTRYPOINT ./heroku/start.sh