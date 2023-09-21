export GOBIN=$(PWD)/bin
export PROTOBUF_ROOT=$(PWD)/_vendor/protobuf-21.12

.PHONY: install test gen-conformance gen-include genall

install:
	go install -tags protolegacy google.golang.org/protobuf/cmd/protoc-gen-go
	go install -tags protolegacy ./cmd/protoc-gen-go-vtproto
# 	go install -tags protolegacy github.com/gogo/protobuf/protoc-gen-gofast

gen-conformance:
	$(PROTOBUF_ROOT)/src/protoc \
		--proto_path=$(PROTOBUF_ROOT) \
		--go_out=conformance --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		--go-vtproto_out=conformance --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
		-I$(PROTOBUF_ROOT)/src \
		--go_opt=Msrc/google/protobuf/test_messages_proto2.proto=internal/conformance \
		--go_opt=Msrc/google/protobuf/test_messages_proto3.proto=internal/conformance \
		--go_opt=Mconformance/conformance.proto=internal/conformance \
		--go-vtproto_opt=Msrc/google/protobuf/test_messages_proto2.proto=internal/conformance \
		--go-vtproto_opt=Msrc/google/protobuf/test_messages_proto3.proto=internal/conformance \
		--go-vtproto_opt=Mconformance/conformance.proto=internal/conformance \
		src/google/protobuf/test_messages_proto2.proto \
		src/google/protobuf/test_messages_proto3.proto \
		conformance/conformance.proto

gen-include:
	$(PROTOBUF_ROOT)/src/protoc \
		--proto_path=include \
		--go_out=include --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		-I$(PROTOBUF_ROOT)/src \
		github.com/planetscale/vtprotobuf/vtproto/ext.proto
	mv include/github.com/planetscale/vtprotobuf/vtproto/*.go ./vtproto

gen-testproto:
	$(PROTOBUF_ROOT)/src/protoc \
		--proto_path=testproto \
		--proto_path=include \
		--go_out=. --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		--go-vtproto_out=allow-empty=true:. --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
		-I$(PROTOBUF_ROOT)/src \
		testproto/empty/empty.proto \
		testproto/pool/pool.proto \
		testproto/pool/pool_with_slice_reuse.proto \
		testproto/pool/pool_with_oneof.proto \
		testproto/proto3opt/opt.proto \
		testproto/proto2/scalars.proto \
		testproto/unsafe/unsafe.proto \
		|| exit 1;

genall: install gen-include gen-conformance gen-testproto

test: install gen-conformance
	go test -short ./...
	go test -count=1 ./conformance/...
	GOGC="off" go test -count=1 ./testproto/pool/...
