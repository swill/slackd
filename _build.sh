#!/usr/bin/env bash

goxc -d=./bin -bc="linux,!arm darwin"
rm -rf debian