#!/bin/sh
# Copyright 2020 Google LLC
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

# turn on bash's job control
set -e

# Start the ext_authz service
GRPC_PORT=4000 ./custom-plugin &

# Start apigee remote-service
./apigee-remote-service-envoy &

# Start envoyproxy
./usr/local/bin/envoy -c /etc/envoy/envoy.yaml #--component-log-level filter:trace,ext_authz:trace

#fg %1