#!/usr/bin/env bash

protoc --go_out=plugins=grpc:. im_protobuf.proto
