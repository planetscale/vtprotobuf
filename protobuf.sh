#/bin/bash
set -e

if [ "$#" -ne 1 ]; then
    echo "need to provide the directory to use" >&2
    exit 2
fi

DEST=$1

curl -sS -L -o "$DEST/protobuf-all-3.16.0.tar.gz" http://github.com/protocolbuffers/protobuf/releases/download/v3.16.0/protobuf-all-3.16.0.tar.gz

cd "$DEST"

tar zxf protobuf-all-3.16.0.tar.gz
cd protobuf-3.16.0
./configure --quiet

make

cd conformance/
make

echo "Dowloaded and compiled protobuf 3.16.0 to $DEST"

exit 0