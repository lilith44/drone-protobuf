#!/usr/bin/env bash
docker build \
--build-arg ALPINE_VERSION="${ALPINE_VERSION:-"3.12"}" \
--build-arg GO_VERSION="${GO_VERSION:-"1.15.3"}" \
--build-arg GRPC_GATEWAY_VERSION="${GRPC_GATEWAY_VERSION:-"2.0.1"}" \
--build-arg GRPC_VERSION="${GRPC_VERSION:-"1.33.1"}" \
--build-arg PROTOC_GEN_DOC_VERSION="${PROTOC_GEN_DOC_VERSION:-"1.3.2-ttn"}" \
--build-arg PROTOC_GEN_FIELDMASK_VERSION="${PROTOC_GEN_FIELDMASK_VERSION:-"0.4.5"}" \
--build-arg PROTOC_GEN_GO_VERSION="${PROTOC_GEN_GO_VERSION:-"1.4.3"}" \
--build-arg PROTOC_GEN_GOGO_VERSION="${PROTOC_GEN_GOGO_VERSION:-"1.3.1"}" \
--build-arg PROTOC_GEN_GOGOTTN_VERSION="${PROTOC_GEN_GOGOTTN_VERSION:-"3.0.14"}" \
--build-arg PROTOC_GEN_HUGODATA_VERSION="${PROTOC_GEN_HUGODATA_VERSION:-"1.0.0"}" \
--build-arg PROTOC_GEN_VALIDATE_VERSION="${PROTOC_GEN_VALIDATE_VERSION:-"0.4.1"}" \
--build-arg UPX_VERSION="${UPX_VERSION:-"3.96"}" \
${@} .
