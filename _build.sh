#!/usr/bin/env bash

gox -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="linux darwin"