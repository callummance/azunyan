# Azunyan
A song request and queueing system for karaoke, written for ICAS.

## Getting started
Dependencies:
- npm
- Golang
- dep

Installation:

```
cd static\frontend
npm install
cd ..\..
dep ensure
go build
```
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
`dbaddr`, `dbaddr`, and `dbcollection` have now been moved to environmental variables to allow deployment to Heroku. These can be included in a `.env` file.

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
git push heroku master
```