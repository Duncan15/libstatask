#!/usr/bin/env bash
cd /Users/cwc/go/src/libstatask
module=libstatask
count=`ps ax|grep $module|grep -v grep|wc -l`
pid=`ps ax|grep $module|grep -v grep|awk '{print $1}'`
if [  $count -gt 0 ]; then
    echo "stop process $module ,pid $pid"
    kill $pid
    sleep 5s
    count=`ps ax|grep $module|grep -v grep|wc -l`
    pid=`ps ax|grep $module|grep -v grep|awk '{print $1}'`
    if [ $count -gt 0 ]; then
        echo "the process still exist, abort!"
        exit 1
    fi
fi

go build
nohup ./$module -m online >nohup.out 2>&1 &

count=`ps ax|grep $module|grep -v grep|wc -l`
pid=`ps ax|grep $module|grep -v grep|awk '{print $1}'`
if [ $count -eq 0 ]; then
    echo "the process fail to startup, abort!"
    exit 1
fi
echo "$module $pid startup"



