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

readonly MY_PATH=$(cd $(dirname "$0") && pwd)
cd $MY_PATH/../src
go test ./...

# Commandline tests.
set -e  # Fail on any error.

# Should show error if no directories are passed.
go run . 2>/dev/null && exit 1
# Should terminate normally.
go run . . 2>/dev/null

# Match expected outputs.
go run . 2>&1 | grep -q "You must specify at least one directory"
go run . --version | grep -q 'Git commit hash'
go run . --help | grep -q 'Usage:'
