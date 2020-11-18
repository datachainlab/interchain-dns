#!/usr/bin/env bash

set -eo pipefail

# Get the path of the cosmos-sdk repo from go/pkg/mod
proto_dirs=$(find . -path ./third_party -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  protoc \
  -I "proto" \
  -I "third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  $(find "${dir}" -name '*.proto')

  # command to generate gRPC gateway (*.pb.gw.go in respective modules) files
  protoc \
  -I "proto" \
  -I "third_party/proto" \
  --grpc-gateway_out=logtostderr=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done

cp -r github.com/datachainlab/cosmos-sdk-interchain-dns/* ./
rm -rf github.com
