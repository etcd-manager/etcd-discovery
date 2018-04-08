#!/usr/bin/env bash

pushd $GOPATH/src/github.com/etcd-manager/etcd-discovery/hack/gendocs
go run main.go
popd
