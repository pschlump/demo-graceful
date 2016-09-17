#!/bin/bash

NN="this-one"

mkdir -p ./out
rm -f ./demo-graceful
go build 2>&1 | color-cat -c red

xx=$( ps -ef | grep demo-graceful | grep -v grep | grep "$NN" | awk '{print $2}' )
if [ "X$xx" == "X" ] ; then	
	:
else
	kill $xx
fi
if [ -x ./demo-graceful ] ; then
	./demo-graceful -c ./make-test.json --this-one &
	SVR_PID=$!
	echo "$SVR_PID" >out/pid-of-svr
	sleep 1
	echo "wait 1 sec for server to start up"
else
	exit 1
fi


