#!/usr/bin/env bash

gox -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="linux/amd64"