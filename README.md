# Azunyan
A song request and queueing system for karaoke, written for ICAS.

## Quick run
If you just need to get a prebuilt copy of azunyan up and running as fast as possible, the Docker Hub image is going to be the way to go. For development work, follow the development guide down below.

First install Docker and latest Docker-Compose:
```bash
sudo apt install docker.io
sudo curl -L https://github.com/docker/compose/releases/download/1.25.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
```
Clone the project
```
git clone https://github.com/callummance/azunyan
```

After ensuring you have docker installed and set up, you will first need to fetch the docker-compose file with the following command:
```bash
curl https://raw.githubusercontent.com/callummance/azunyan/master/docker-compose-prod.yml -o docker-compose.yml
```
At this point you will need to create an `azunyan.conf` and a `ssh_pass.conf` file. The first contains config for the karaoke server itself, wheras the latter just contains the password which may be used to ssh into the docker network (useful for adding new songs). You can fetch examples of both by executing:
```bash
curl -O https://raw.githubusercontent.com/callummance/azunyan/master/ssh_pass.conf
curl -O https://raw.githubusercontent.com/callummance/azunyan/master/azunyan.conf
```
You will, however, want to change the contents of both of these for security reasons.

Finally, to start the queue system simply run
```bash
docker-compose up
```

On a browser go to \<server ip address\>:8080

Note: By default the docker container will start the server on port 8080. If deploying the container on an Amazon EC2 instance, remember to add a custom TCP rule in the security group associated with the EC2 instance.
***
## Development Guide

### Dependencies:
These are required to run the app:
- npm
- Golang
- dep
- mongoDB

### Rough installation notes for Ubuntu 16.04

Follow the below steps if you want to manually build the project for development. If you wish to just run Azunyan, then the Quick Run section detailed above is recommended as it will save you a lot of setup effort.

```bash
sudo apt install golang-1.9 # This actually installs the go binary to /usr/lib for some reason so we will need to do a symlink
sudo ln -s /usr/lib/go-1.9/bin/go /usr/bin/go

# Add $GOPATH to your PATH. This will typically be ~/go
echo "export GOPATH=~/go" >> ~/.bashrc

# Set up Go paths
mkdir -p ~/go/src
mkdir -p ~/go/bin
cd ~/go/src

# Clones the Azunyan project to your GOPATH directory.
go get github.com/callummance/azunyan

# Install dep. This is used to install the backend server dependencies.
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# ^ Above will install dep to $GOPATH/bin/dep. Add $GOPATH/bin to your path

# Install Node and NPM. I recommend doing this with NVM. Saves a lot of headache.
# See https://github.com/nvm-sh/nvm

# Setup a MongoDB database
# See https://www.mongodb.com/try/download/community and https://www.freecodecamp.org/news/learn-mongodb-a4ce205e7739/ for some guides on how to to do this

# Build the app
cd $GOPATH/github.com/callummance/azunyan
npm install
dep ensure
go build
```

This app runs on a MongoDB server. You will need to create an instance of this, either locally or using the Cloud Atlas platform (remember to whitelist your IP address if using the Cloud Atlas platform).
You should then set the corresponding environmental variables. See section Configuration for more info.

Run:

On Linux:

```
./azunyan
```

On Windows:

```
.\azunyan.exe
```

### Configuration
Configuration for this program is stored within `azunyan.conf`. The default file is populated with all fields supported. 
The fields `dbaddr`, `dbaddr`, and `dbcollection` can also be retrieved from environmental variables to allow deployment to Heroku. These can be included in a `.env` file. An example `.env` file can be seen in `.env-example` (Note in MongoDB, a table is known as a 'collection')

### Heroku
It is also possible to deploy the docker container to Heroku. The advantages to this is you can get your own free heroku domain name
and no need to specify port 8080. Furthermore, you don't need your own server which altogether means it'll be free to run.
However, the cons is accessing the MongoDB is much slower (1-3 second delays) due to the fact that Heroku cannot install the MongoDB 
on the same server that is hosting the application. This means the database is remote to the app so requires
communication via internet. As such, we do not recommend this strategy if a lot of people are expected to be requesting songs at the
same time.
To deploy to Heroku, first install the Heroku cli:
```
curl https://cli-assets.heroku.com/install.sh | sh
```
Install Docker if you haven't already. Then run
```
heroku login
heroku create <app name>
heroku/deploy_to_heroku.sh
heroku ps:scale web=1
```

### Useful hints when debugging:
Remember to do a hard refresh in the browser to see changes in JS files.
`docker system prune -af` can help if you want to completely rebuild the docker containers.
