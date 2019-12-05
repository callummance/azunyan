#!/bin/sh
mongo $MONGODB_URI ./container-scripts/mongo-init.js
/go/bin/azunyan