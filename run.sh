#!/bin/sh

port=8080

if [[ x$1 != x ]]
then
    port=$1
fi

while true;
do
    /src/shiori serve -p $port
    sleep 5
done
