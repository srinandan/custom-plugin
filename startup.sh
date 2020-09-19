#!/usr/bin/env sh

# turn on bash's job control
set -e

# Start the first process
./custom-plugin &

./usr/local/bin/envoy -c /etc/envoy/envoy.yaml

