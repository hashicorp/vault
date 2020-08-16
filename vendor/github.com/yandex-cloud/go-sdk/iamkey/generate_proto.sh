#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
protoc  \
    --proto_path ../../public-api/ \
    --proto_path . \
    --go_out=Myandex/cloud/iam/v1/key.proto=bb.yandex-team.ru/cloud/cloud-go/genproto/publicapi/yandex/cloud/iam/v1:$GOPATH/src *.proto

