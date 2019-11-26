# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

proto:
	protoc --go_out=paths=source_relative:. types/types.proto
