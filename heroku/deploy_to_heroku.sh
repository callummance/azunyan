#!/bin/bash
heroku stack:set container
MONGO_INSTALLED=$(heroku addons | grep mongo)
if [ -z "$MONGO_INSTALLED" ]
then
    heroku addons:create mongolab:sandbox
fi
git push -f heroku master
DBADDR=$(heroku config:get MONGODB_URI)
DBNAME=$(basename $DBADDR)
heroku config:set dbaddr=$DBADDR dbname=$DBNAME dbcollection=songs
