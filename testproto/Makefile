export GOBIN=$(PWD)/../bin
export PROTOBUF_ROOT=$(HOME)/src/protobuf-3.16.0
VTROOT=$(HOME)/src/vitess

.PHONY: regenerate

regenerate:
	$(VTROOT)/bin/protoc \
		--proto_path=. \
		--go_out=. --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		--go-vtproto_out=. --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
		-I$(PROTOBUF_ROOT)/src/google/protobuf \
		-I$(PWD)/../include \
		pool.proto