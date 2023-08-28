#!/usr/bin/env bash

gox -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="linux/amd64" -osarch="darwin/arm64" -osarch="darwin/amd64" -osarch="windows/amd64"