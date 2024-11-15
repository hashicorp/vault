# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

proto:
	buf generate
	
proto-lint:
	buf lint

proto-format:
	buf format -w
