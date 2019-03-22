# Azunyan
A song request and queueing system for karaoke, written for ICAS.

## Getting started
Run:

```
cd static\frontend
npm install
cd ..\..
go build
```


### Configuration
Configuration for this program is stored within `azunyan.conf`. The default file is populated with all fields supported. 

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
