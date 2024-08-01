#!/bin/bash

ginWatch() {
  gin -a $1 -i -p $2 run .
}

if which gin; then
  ginWatch $1 $(($1 + 1))
else
  go install github.com/codegangsta/gin@latest
  ginWatch $1 $(($1 + 1))
fi
