#!/bin/bash

export MONGODB_CONNECTIONSTRING="mongodb://spyglass:spyglass@localhost:27017/spyglass?authSource=admin"
./spyglass server -v DEBUG

