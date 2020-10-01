##             MonitorAPIs microservice Makefile
##
##  Simple Makefile containing implementation of targets for generating protobuf file
##
##  To generate swagger, angular models and go messages&grpc use this command:
##    $ make gen
SHELL=/bin/bash

PROTO_INCLUDES=-Iproto -I3rdparty -I3rdparty/protoc-gen-validate -I3rdparty/api-common-protos -I=.
PROTO_GRPC_ARGS=paths=source_relative
PROTO_SWAGGER_ARGS=logtostderr=true,json_names_for_fields=true
##
##  \e[1mTargets\e[0m
##   \e[34mhelp\e[0m
##       Shows this help
help:
	@echo -e "$$(sed -n 's/^##//p' Makefile)"

##   \e[34mgen\e[0m
##       Shortcut for generate
gen:
	make generate SERVICE_NAME=notification SERVICE_VERSION=v1

##   \e[34mgenerate\e[0m
##       Generates Go
generate: generate/go

##   \e[34mgenerate/go\e[0m
##       Generates go grpc files and messages from proto file
generate/go:
	protoc $(PROTO_INCLUDES) \
		  proto/${SERVICE_VERSION}/${SERVICE_NAME}.proto \
		   --go_out=$(PROTO_GRPC_ARGS):go-sdk/${SERVICE_NAME} \
		   --go-grpc_out=$(PROTO_GRPC_ARGS):go-sdk/${SERVICE_NAME} \
			--grpc-gateway_out=logtostderr=true,allow_patch_feature=false,paths=source_relative:go-sdk/${SERVICE_NAME} \
		   --validate_out="lang=go,$(PROTO_GRPC_ARGS):go-sdk/${SERVICE_NAME}"

	protoc-go-inject-tag -input=./go-sdk/${SERVICE_NAME}/${SERVICE_VERSION}/${SERVICE_NAME}.pb.go