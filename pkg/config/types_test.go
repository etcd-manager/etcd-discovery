package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/appscode/kutil/meta"
)

/*
/opt/etcd-v3.2.18-linux-amd64/etcd --name infra1 \
  --initial-advertise-peer-urls https://127.0.0.1:2380 \
  --listen-peer-urls https://127.0.0.1:2380 \
  --listen-client-urls https://127.0.0.1:2379 \
  --advertise-client-urls https://127.0.0.1:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra1=https://127.0.0.1:2380 \
  --initial-cluster-state new \
  --data-dir=/tmp/infra1 \
  --cert-file=etcd.local.config/certificates/db-server.crt \
  --key-file=etcd.local.config/certificates/db-server.key \
  --trusted-ca-file=etcd.local.config/certificates/db-ca.crt \
  --client-cert-auth=true \
  --peer-cert-file=etcd.local.config/certificates/peer-localhost.crt \
  --peer-key-file=etcd.local.config/certificates/peer-localhost.key \
  --peer-trusted-ca-file=etcd.local.config/certificates/peer-ca.crt \
  --peer-client-cert-auth=true
*/
func TestEtcdFlags(t *testing.T) {
	f := NewEtcdFlags()

	f.Name = "infra1"
	f.InitialAdvertisePeerURLs.Insert("127.0.0.1")
	f.ListenPeerURLs.Insert("127.0.0.1")
	f.ListenClientURLs.Insert("127.0.0.1")
	f.AdvertiseClientURLs.Insert("127.0.0.1")
	f.InitialClusterToken = "etcd-cluster-1"
	f.InitialCluster.Insert("infra1", "127.0.0.1")
	f.InitialClusterState = "new"
	f.DataDir = "/tmp/infra1"
	f.CertFile = "etcd.local.config/certificates/db-server.crt"
	f.KeyFile = "etcd.local.config/certificates/db-server.key"
	f.TrustedCAFile = "etcd.local.config/certificates/db-ca.crt"
	f.ClientCertAuth = true
	f.PeerCertFile = "etcd.local.config/certificates/peer-localhost.crt"
	f.PeerKeyFile = "etcd.local.config/certificates/peer-localhost.key"
	f.PeerTrustedCAFile = "etcd.local.config/certificates/peer-ca.crt"
	f.PeerClientCertAuth = true

	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(data))

	m := map[string]string{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Error(err)
	}

	args := meta.BuildArgumentListFromMap(m, nil)
	fmt.Println(strings.Join(args, ",\n"))
}
