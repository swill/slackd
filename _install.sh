#!/usr/bin/env bash

git fetch
git rebase origin/master
tar zxfv ./bin/snapshot/slackd_linux_amd64.tar.gz -C . --strip-components 1