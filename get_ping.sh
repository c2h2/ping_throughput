#!/usr/bin/env bash


if [ "$(uname)" == "Darwin" ]; then
   ping -W 3000 -c 1 $1 | head -2 |tail -1 |cut -d" " -f7|cut -d"=" -f2 
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
   ping -W 3000 -c 1 $1 | head -2 |tail -1 |cut -d" " -f8|cut -d"=" -f2 
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
   echo $1
    # Do something under 32 bits Windows NT platform
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW64_NT" ]; then
   echo $1
fi
