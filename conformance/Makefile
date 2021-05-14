export GOBIN=$(PWD)/bin
export PROTOBUF_ROOT=$(HOME)/src/protobuf-3.16.0
VTROOT=$(HOME)/src/vitess

.PHONY: install regenerate test

install:
	go install -tags protolegacy google.golang.org/protobuf/cmd/protoc-gen-go
	go install -tags protolegacy ../cmd/protoc-gen-go-vtproto
	go install -tags protolegacy github.com/gogo/protobuf/protoc-gen-gofast

regenerate-gogo: install
	$(VTROOT)/bin/protoc \
		--proto_path=$(PROTOBUF_ROOT) \
		--gofast_out=. --plugin protoc-gen-gofast="${GOBIN}/protoc-gen-gofast" \
		-I$(PROTOBUF_ROOT)/src/google/protobuf \
		--gofast_opt=Msrc/google/protobuf/test_messages_proto3.proto=internal/conformance_gogo \
		src/google/protobuf/test_messages_proto3.proto

regenerate: install
	$(VTROOT)/bin/protoc \
		--proto_path=$(PROTOBUF_ROOT) \
		--go_out=. --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		--go-vtproto_out=. --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
		-I$(PROTOBUF_ROOT)/src/google/protobuf \
		--go_opt=Msrc/google/protobuf/test_messages_proto2.proto=internal/conformance \
		--go_opt=Msrc/google/protobuf/test_messages_proto3.proto=internal/conformance \
		--go_opt=Mconformance/conformance.proto=internal/conformance \
		--go-vtproto_opt=Msrc/google/protobuf/test_messages_proto2.proto=internal/conformance \
		--go-vtproto_opt=Msrc/google/protobuf/test_messages_proto3.proto=internal/conformance \
		--go-vtproto_opt=Mconformance/conformance.proto=internal/conformance \
		src/google/protobuf/test_messages_proto2.proto \
		src/google/protobuf/test_messages_proto3.proto \
		conformance/conformance.proto

test: regenerate
	go test -count=1 ./...