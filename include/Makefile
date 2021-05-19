export GOBIN=$(PWD)/../bin
export PROTOBUF_ROOT=$(HOME)/src/protobuf-3.16.0
VTROOT=$(HOME)/src/vitess

.PHONY: regenerate

regenerate:
	$(VTROOT)/bin/protoc \
		--proto_path=. \
		--go_out=. --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		-I$(PROTOBUF_ROOT)/src/google/protobuf \
		vitess.io/vtprotobuf/vtproto/ext.proto
	mv vitess.io/vtprotobuf/vtproto/*.go $(PWD)/../vtproto