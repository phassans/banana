#!/bin/sh
#compile go build
go build .

#check if banana process is running
ps -ef | grep ./banana
if [ $? -eq 0 ]
then
  echo "banana Running..."
  echo "killing process..."
  ps -ef | grep ./banana | grep -v grep | awk '{print $2}' | xargs kill
  if [ $? -eq 0 ]
    then
    echo "banana process killed!"
  else
    echo "could not kill banana process" >&2
    exit 1
  fi
else
  echo "process not running"
fi

nohup ./banana &
if [ $? -eq 0 ]
then
  echo "banana process restarted"
  tail -f nohup.out
else
  echo "failed to start banana process" >&2
fi
