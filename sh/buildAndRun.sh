#!/bin/bash
PROJECTNAME="summer"
CONFIGNAME="config_dev.conf"


PROJECT_PATH=$(cd "$(dirname "$0")";cd ../; pwd)
BINPATH=$PROJECT_PATH"/bin/"
GOCMD=$(which go)
#
#build
cd $PROJECT_PATH
$GOCMD build  -o $BINPATH .

# kill

pid=`ps -ef | grep $PROJECTNAME'go' | grep -v grep | awk '{print $2}'`

if [ "$pid" ]
then
  echo '杀死进程:'$pid
  kill -9 $pid
fi

#run
cd $BINPATH
rm -rf ./run.log
nohup ./$PROJECTNAME -config ./$CONFIGNAME killtag=$PROJECTNAME"go"  >run.log 2>&1 &



