#!/bin/bash

heroku stack:set container
DBADDR=$(heroku config:get dbaddr)
DBNAME=$(heroku config:get dbname)

while (( "$#" )); do
    case "$1" in
    -dbaddr)
        DBADDR="$2"
        ;;
    -dbname)
        DBNAME="$2"
        ;;
    esac
done
if [ -z "$DBADDR" ];
then
    echo "No mongoDB URI found. Please specify the mongoDB URI with the -dbaddr option"
    exit
fi
if [ -z "$DBNAME" ];
then
    echo "No mongoDB database name found. Please specify the mongoDB's database name with the -dbname option"
    exit
fi

git push -f heroku master
heroku config:set dbaddr=$DBADDR dbname=$DBNAME dbcollection=songs
