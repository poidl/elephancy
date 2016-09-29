#!/bin/bash
# fingerprint file name with md5 hash and Date-Modified

if ! [ -f $1 ]; then
  echo "File not found: $1"
  exit 1
fi

if [ $# -ne 2 ]; then
  echo 'Usage: fingerprint FILENAME DATESTRING'
  echo 'DATESTRING e.g. of the form 20150122'
  exit 1
fi

filename=$(basename "$1")
suffix=${filename##*.}
dirname=$(dirname "$1")

# check if the filename is of form "name_DATE_MD5HASH.suffix"
if [[ $1 =~ ^[^_]*_[0-9]*_[0-9a-f]*\.[a-zA-Z0-9]*$ ]]; then
  prefix=${filename%%_*}
# check if the filename is of form name.suffix
elif [[ $1 =~ ^[a-zA-Z0-9]*\.[a-zA-Z0-9]*$ ]]; then
  prefix="${filename%.*}"
else
  echo "Unexpected file name"
  exit 1
fi

hash=$(md5sum $1|awk '{print $name}')
hash=$(echo $hash|cut -c1-7) # keep first 7 characters

mv $1 $dirname"/"$prefix"_"$2"_"$hash"."$suffix
