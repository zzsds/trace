PWD:=$(shell pwd)
PROTO=match

.PHONY: proto
proto:
	@echo execute ${PROTO} proto file generate
	protoc --proto_path=. --go_out=. proto/${PROTO}/${PROTO}.proto
