# Azunyan
A song request and queueing system for karaoke, written for ICAS.

## Getting started
Dependencies:
- npm
- Golang
- dep

## Installation:
Follow the below steps if you want to manually build the project for development. If you wish to just run Azunyan, then the docker-compose section detailed lower down is recommended as it will save you a lot of setup effort.

### Rough installation notes for Ubuntu 16.04

```bash
sudo apt install golang-1.9 # This actually installs the go binary to /usr/lib for some reason so we will need to do a symlink
sudo ln -s /usr/lib/go-1.9/bin/go /usr/bin/go

# Add $GOPATH to your PATH. This will typically be ~/go
echo "export GOPATH=~/go" >> ~/.bashrc

# Set up Go paths
mkdir -p ~/go/src
mkdir -p ~/go/bin
cd ~/go/src

# Install Azunyan
go get github.com/callummance/azunyan

# Install dep
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# ^ Above will install dep to $GOPATH/bin/dep. Add $GOPATH/bin to your path

# Install Node and NPM. I recommend doing this with NVM. Saves a lot of headache.
# See https://github.com/nvm-sh/nvm

cd $GOPATH/github.com/callummance/azunyan
npm install
dep ensure
go build
```

This app runs on a MongoDB server. You will need to create an instance of this, either locally or using the Cloud Atlas platform.
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
`dbaddr`, `dbaddr`, and `dbcollection` have now been moved to environmental variables to allow deployment to Heroku. These can be included in a `.env` file. An example `.env` file can be seen in `.env-example` (Note in MongoDB, a table is known as a 'collection')

### Docker
First, you need to modify the configuration file as described above.


Install Docker and latest Docker-Compose:
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
Build the image by running 
```
docker-compose build
```
and the image can then be started by running
```
docker-compose up
```
On a browser go to \<server ip address\>:8080

Note: By default the docker container will start the server on port 8080. If deploying the container on an Amazon EC2 instance, remember to add a custom TCP rule in the security group associated with the EC2 instance.

### Heroku
To deploy to Heroku, first install the Heroku cli:
```
curl https://cli-assets.heroku.com/install.sh | sh
```

 then run
```
heroku create
heroku buildpacks:add --index 2 heroku/nodejs
git push heroku master
```