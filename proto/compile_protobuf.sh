#!/usr/bin/env bash
set -euo pipefail

protoc *.proto --go_out=:.