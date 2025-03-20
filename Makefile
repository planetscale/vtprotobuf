export GOBIN=$(PWD)/bin
export PROTOBUF_ROOT=$(PWD)/_vendor/protobuf-30.1
export PROTOC_PATH=$(PROTOBUF_ROOT)/protoc

.PHONY: install test gen-conformance gen-include gen-wkt genall bin/protoc-gen-go bin/protoc-gen-go-vtproto

install: bin/protoc-gen-go-vtproto bin/protoc-gen-go

bin/protoc-gen-go-vtproto:
	go install -buildvcs=false -tags protolegacy ./cmd/protoc-gen-go-vtproto

bin/protoc-gen-go:
	go install -tags protolegacy google.golang.org/protobuf/cmd/protoc-gen-go@latest

gen-conformance: install
	$(PROTOC_PATH) \
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
	$(PROTOC_PATH) \
		--proto_path=include \
		--go_out=include --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		-I$(PROTOBUF_ROOT)/src \
		github.com/planetscale/vtprotobuf/vtproto/ext.proto
	mv include/github.com/planetscale/vtprotobuf/vtproto/*.go ./vtproto

gen-wkt: bin/protoc-gen-go-vtproto
	$(PROTOC_PATH) \
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
	$(PROTOC_PATH) \
		--proto_path=testproto \
		--proto_path=include \
		--go_out=. --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		--go-vtproto_out=allow-empty=true:. --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
		-I$(PROTOBUF_ROOT)/src \
		testproto/ignore_unknown_fields/opt.proto \
		testproto/empty/empty.proto \
		testproto/pool/pool.proto \
		testproto/pool/pool_with_slice_reuse.proto \
		testproto/pool/pool_with_oneof.proto \
		testproto/proto3opt/opt.proto \
		testproto/proto2/scalars.proto \
		testproto/unsafe/unsafe.proto \
		testproto/unique/unique.proto \
		|| exit 1;
	$(PROTOC_PATH) \
		--proto_path=testproto \
		--proto_path=include \
		--go_out=. --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		--go-vtproto_opt=paths=source_relative \
		--go-vtproto_opt=buildTag=vtprotobuf \
		--go-vtproto_out=allow-empty=true:./testproto/buildtag --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
		-I$(PROTOBUF_ROOT)/src \
		testproto/empty/empty.proto \
		|| exit 1;

get-grpc-testproto: install
	$(PROTOC_PATH) \
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
	$(PROTOC_PATH) \
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
	go test -count=1 ./testproto/ignore_unknown_fields/...
	GOGC="off" go test -count=1 ./testproto/pool/...
