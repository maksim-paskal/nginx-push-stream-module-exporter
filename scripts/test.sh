#!/usr/bin/env bash

# Copyright paskal.maksim@gmail.com
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

export CGO_ENABLED=0
export GOFLAGS="-trimpath"
export MY_POD_NAMESPACE=default

set -ex

rm -rf main
go build -ldflags "-X main.buildTime=$(date +"%Y%m%d%H%M%S")" -o main -v ./cmd/
./main -log.level=DEBUG -log.pretty -nginx.address=http://127.0.0.1:18102 $*