# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

proto:
	protoc types.proto --go_out=paths=source_relative:.
	sed -i -e 's/Iv/IV/' -e 's/Hmac/HMAC/' types.pb.go

.PHONY: proto
