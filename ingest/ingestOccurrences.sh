#!/bin/sh

if [ -z "$1" ]; then
  echo "usage: $0 <occurrence data>"
  exit 1
fi

mongoimport --db iroiro --collection occurrences --file "$1" \
  --type=tsv --headerline

