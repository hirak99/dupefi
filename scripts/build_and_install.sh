#!/bin/bash
# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -uexo pipefail
readonly MY_PATH=$(cd $(dirname "$0") && pwd)

$MY_PATH/run_tests.sh

cd $MY_PATH/../src

LOCAL_MOD=""
if ! (git diff --exit-code > /dev/null) \
  || ! [[ -z "$(git status --porcelain)" ]] ; then
  LOCAL_MOD="_locally_modified"
fi

GITHASH=$(git rev-parse HEAD)$LOCAL_MOD
echo $GITHASH

GITDATE=$(date --date=@$(git log -1 --format="%at") +%Y_%m_%d)

# To get variable names for X, try go tool nm /usr/bin/dupefi
BUILDINFO="nomen_aliud/dupefi/buildinfo"
go build -ldflags "-X '$BUILDINFO.Githash=$GITHASH' -X '$BUILDINFO.BuildTime=$GITDATE'"
# Or -
# go build -ldflags "-X 'main.Githash=$GITHASH'"

sudo install -Dm755 ./dupefi /usr/bin/dupefi
rm dupefi
