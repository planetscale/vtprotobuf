export GOBIN=$(PWD)/bin
export PROTOBUF_ROOT=$(PWD)/_vendor/protobuf-21.12

.PHONY: install test gen-conformance gen-include gen-wkt genall bin/protoc-gen-go bin/protoc-gen-go-vtproto

install: bin/protoc-gen-go-vtproto bin/protoc-gen-go

bin/protoc-gen-go-vtproto:
	go install -tags protolegacy ./cmd/protoc-gen-go-vtproto

bin/protoc-gen-go:
	go install -tags protolegacy google.golang.org/protobuf/cmd/protoc-gen-go

gen-conformance: install
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

gen-include: bin/protoc-gen-go
	$(PROTOBUF_ROOT)/src/protoc \
		--proto_path=include \
		--go_out=include --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		-I$(PROTOBUF_ROOT)/src \
		github.com/planetscale/vtprotobuf/vtproto/ext.proto
	mv include/github.com/planetscale/vtprotobuf/vtproto/*.go ./vtproto

gen-wkt: bin/protoc-gen-go-vtproto
	$(PROTOBUF_ROOT)/src/protoc \
		-I$(PROTOBUF_ROOT)/src \
		--plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
		--go-vtproto_out=. \
		--go-vtproto_opt=module=google.golang.org/protobuf,wrap=true \
		$(PROTOBUF_ROOT)/src/google/protobuf/any.proto \
        $(PROTOBUF_ROOT)/src/google/protobuf/duration.proto \
        $(PROTOBUF_ROOT)/src/google/protobuf/empty.proto \
        $(PROTOBUF_ROOT)/src/google/protobuf/field_mask.proto \
        $(PROTOBUF_ROOT)/src/google/protobuf/timestamp.proto \
        $(PROTOBUF_ROOT)/src/google/protobuf/wrappers.proto \
        $(PROTOBUF_ROOT)/src/google/protobuf/struct.proto

gen-testproto: get-grpc-testproto gen-wkt-testproto install
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

get-grpc-testproto: install
	$(PROTOBUF_ROOT)/src/protoc \
		--proto_path=. \
		--proto_path=include \
		--go_out=. --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		--go-vtproto_out=. --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
		-I$(PROTOBUF_ROOT)/src \
		-I. \
		--go_opt=paths=source_relative \
		--go_opt=Mtestproto/grpc/inner/inner.proto=github.com/planetscale/vtprotobuf/testproto/grpc/inner \
		--go-vtproto_opt=paths=source_relative \
        --go-vtproto_opt=Mtestproto/grpc/inner/inner.proto=github.com/planetscale/vtprotobuf/testproto/grpc/inner \
		testproto/grpc/inner/inner.proto \
		testproto/grpc/grpc.proto \
		|| exit 1;

gen-wkt-testproto: install
	$(PROTOBUF_ROOT)/src/protoc \
    	--proto_path=testproto \
    	--proto_path=include \
    	--go_out=. --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
    	--go-vtproto_out=allow-empty=true:. --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
    	-I$(PROTOBUF_ROOT)/src \
    	testproto/wkt/wkt.proto \
    	|| exit 1;

genall: gen-include gen-conformance gen-testproto gen-wkt

test: install gen-conformance
	go test -short ./...
	go test -count=1 ./conformance/...
	GOGC="off" go test -count=1 ./testproto/pool/...
