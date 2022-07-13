#!/bin/bash
readonly MY_PATH=$(cd $(dirname "$0") && pwd)
cd $MY_PATH/../src
go test .
