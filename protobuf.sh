#!/bin/bash

set -e

PROTOBUF_VERSION=30.1
ROOT=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
PROTOBUF_PATH="${ROOT}/_vendor/protobuf-${PROTOBUF_VERSION}"

if [ -f "$PROTOBUF_PATH/protoc" ]; then
    echo "protoc found in $PROTOBUF_PATH"
    exit 0
fi

mkdir -p _vendor
curl -sS -L -o "$ROOT/_vendor/pb.tar.gz" http://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protobuf-${PROTOBUF_VERSION}.tar.gz

cd "$ROOT/_vendor"
tar zxf pb.tar.gz

cd protobuf-${PROTOBUF_VERSION}
cmake . -Dprotobuf_BUILD_CONFORMANCE=ON -DCMAKE_CXX_STANDARD=17 -Dprotobuf_FORCE_FETCH_DEPENDENCIES=ON
make protoc conformance_test_runner -j$(nproc)

echo "Dowloaded and compiled protobuf $PROTOBUF_VERSION to $PROTOBUF_PATH"
ls -l src/
ls -l conformance/
