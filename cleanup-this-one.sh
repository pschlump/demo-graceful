#!/bin/bash

NN="this-one"

rm -f ./svr
go build 2>&1 | color-cat -c red

xx=$( ps -ef | grep svr | grep -v grep | grep "$NN" | awk '{print $2}' )
if [ "X$xx" == "X" ] ; then	
	:
else
	kill $xx
fi
if [ -x ./svr ] ; then
	./svr --this-one &
	SVR_PID=$!
	echo "$SVR_PID" >pid-of-svr
	sleep 1
	echo "wait 1 sec for server to start up"
else
	exit 1
fi


