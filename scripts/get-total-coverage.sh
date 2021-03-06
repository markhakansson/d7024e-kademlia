#!/bin/bash

go test -v -coverpkg=./... -coverprofile=profile.cov ./... >/dev/null 2>&1
go tool cover -func profile.cov | tail -n 1 | awk '{print $3}' | grep -o -P '[0-9]*(\.[0-9]*)'
