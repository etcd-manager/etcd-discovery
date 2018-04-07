#!/bin/bash

curl -k -vv \
  -XPOST https://127.0.0.1:2381/apis/discovery.etcd-manager.com/v1alpha1/pings \
  -H "Content-Type: application/json" \
  -d '{"apiVersion": "discovery.etcd-manager.com/v1alpha1", "kind": "Ping"}'
