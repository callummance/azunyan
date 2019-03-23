# Azunyan
A song request and queueing system for karaoke, written for ICAS.

## Getting started
Dependencies:
- npm
- Golang
- dep

Installation:

```
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


Build the image by running 
```
docker-compose build
```
and the image can then be started by running
```
docker-compose up
```

### Heroku
To deploy to Heroku, first install the Heroku cli then run
```
heroku create
heroku buildpacks:add --index 2 heroku/nodejs
git push heroku master
```