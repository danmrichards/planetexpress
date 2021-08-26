#!/usr/bin/env bash
set -eu;

inputSpecPath="$1";
inputTarget="$2";
output="$3";

codegenCmd='github.com/deepmap/oapi-codegen/cmd/oapi-codegen';

# This will implicitly fetch and compile the third party
# code generator at the correct version
go run "${codegenCmd}" \
  -o "${output}" \
  -package api \
  -generate "${inputTarget}" \
  "${inputSpecPath}";
