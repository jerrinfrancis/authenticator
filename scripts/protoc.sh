#!/bin/sh
#
protoc --proto_path=proto --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/authenticator/auth.proto
