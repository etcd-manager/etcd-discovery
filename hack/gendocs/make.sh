#!/usr/bin/env bash

pushd $GOPATH/src/github.com/etcd-manager/etcd-discovery/hack/gendocs
go run main.go

cd $GOPATH/src/github.com/etcd-manager/etcd-discovery/docs/reference
sed -i 's/######\ Auto\ generated\ by.*//g' *
popd
