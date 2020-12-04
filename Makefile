PWD:=$(shell pwd)
NAME=match

.PHONY: proto
proto:
	@echo execute ${NAME} proto file generate
	protoc --proto_path=. --go_out=. proto/${NAME}/${NAME}.proto
