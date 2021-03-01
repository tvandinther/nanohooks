#!/usr/bin/env bash
set -eu

protoc *.proto --go_out=:.