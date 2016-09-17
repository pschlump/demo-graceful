#!/bin/bash

if ps -ef | grep -v grep | grep `cat out/pid-of-svr` ; then
	exit 1
fi

if ps -ef | grep -v grep | grep demo-graceful | grep this-one ; then
	exit 1
fi

exit 0
