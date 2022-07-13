#!/bin/bash
set -uexo pipefail
readonly MY_PATH=$(cd $(dirname "$0") && pwd)
cd $MY_PATH/../src
go build .
sudo install -Dm755 ./duphunter /usr/bin/duphunter
rm duphunter
