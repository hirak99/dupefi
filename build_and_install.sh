#!/bin/bash
set -uexo pipefail
cd src
go build .
sudo install -Dm755 ./duphunter /usr/bin/duphunter
rm duphunter
