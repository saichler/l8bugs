#!/usr/bin/env bash

set -e

wget https://raw.githubusercontent.com/saichler/l8types/refs/heads/main/proto/api.proto
wget https://raw.githubusercontent.com/saichler/l8erp/refs/heads/main/proto/erp-common.proto

# Use the protoc image to run protoc.sh and generate the bindings.

# Shared ERP types (must be first - bugs.proto depends on it)
docker run --user "$(id -u):$(id -g)" -e PROTO=erp-common.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest

# L8Bugs types
docker run --user "$(id -u):$(id -g)" -e PROTO=bugs.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest

rm api.proto
rm erp-common.proto

# Now move the generated bindings to the models directory and clean up
rm -rf ../go/types
mkdir -p ../go/types
mv ./types/* ../go/types/.
rm -rf ./types

rm -rf *.rs

cd ../go
find . -name "*.go" -type f -exec sed -i 's|"./types/l8api"|"github.com/saichler/l8types/go/types/l8api"|g' {} +
find . -name "*.go" -type f -exec sed -i 's|"./types/erp"|"github.com/saichler/l8bugs/go/types/erp"|g' {} +
