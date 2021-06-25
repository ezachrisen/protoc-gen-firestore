LOCAL_PATH = $(shell pwd)

.PHONY: example proto install gen-tag test

example: proto install gen-tag
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	-I ./example \
	--firestore_out=module="alticeusa.com/maui/protoc-gen-firestore/example"m:./output example/example.proto

proto:
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	--go_out=:./output example/example.proto --go_opt=module="alticeusa.com/maui/protoc-gen-firestore/example"


install:
	go install .

gen-tag:
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	--go_out=paths=source_relative:. firestore/firestore.proto


test:
	go test ./...
