package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/appscode/go/encoding/json/types"
	"github.com/appscode/kutil/meta"
	"github.com/etcd-manager/etcd-discovery/pkg/etcdclient"
)

type EtcdCluster struct {
	ClusterName     string
	ClusterToken    string
	ClusterSize     int
	BackupStorePath string
	DataDir         string

	InitialClusterState ClusterState
	InitialCluster      map[string]string
}

const (
	ClientPort            = 2379
	PeerPort              = 2380
	DiscoveryPort         = 2381
	QuarantinedClientPort = 8001
)

type EtcdVersion string

// IsV2 returns true if the specified etcdVersion is a 2.x version
func (v EtcdVersion) IsV2() bool {
	return strings.HasPrefix(string(v), "2.")
}

func (v EtcdVersion) GetDockerImage() string {
	return fmt.Sprintf("quay.io/coreos/etcd:%s", v)
}

type EtcdFlags struct {
	Version         EtcdVersion `json:"-"`
	Quarantined     bool        `json:"-"`
	CertificatesDir string      `json:"-"`

	Name                     string        `json:"name"`
	InitialAdvertisePeerURLs *types.URLSet `json:"initial-advertise-peer-urls"`
	ListenPeerURLs           *types.URLSet `json:"listen-peer-urls"`
	ListenClientURLs         *types.URLSet `json:"listen-client-urls"`
	AdvertiseClientURLs      *types.URLSet `json:"advertise-client-urls"`
	InitialClusterToken      string        `json:"initial-cluster-token"`
	InitialCluster           *types.URLMap `json:"initial-cluster"`
	InitialClusterState      string        `json:"initial-cluster-state"`
	DataDir                  string        `json:"data-dir"`
	CertFile                 string        `json:"cert-file"`
	KeyFile                  string        `json:"key-file"`
	TrustedCAFile            string        `json:"trusted-ca-file"`
	ClientCertAuth           types.BoolYo  `json:"client-cert-auth"`
	PeerCertFile             string        `json:"peer-cert-file"`
	PeerKeyFile              string        `json:"peer-key-file"`
	PeerTrustedCAFile        string        `json:"peer-trusted-ca-file"`
	PeerClientCertAuth       types.BoolYo  `json:"peer-client-cert-auth"`
	ForceNewCluster          types.BoolYo  `json:"force-new-cluster"`
	EnableV2                 types.BoolYo  `json:"enable-v2"`
}

func NewEtcdFlags() *EtcdFlags {
	f := &EtcdFlags{
		InitialAdvertisePeerURLs: types.NewURLSet("https", PeerPort),
		ListenPeerURLs:           types.NewURLSet("https", PeerPort),
		ListenClientURLs:         types.NewURLSet("https", ClientPort),
		AdvertiseClientURLs:      types.NewURLSet("https", ClientPort),
		InitialCluster:           types.NewURLMap("https", PeerPort),
	}
	f.ListenPeerURLs.Insert("127.0.0.1")
	f.ListenClientURLs.Insert("127.0.0.1")
	return f
}

func (f *EtcdFlags) ToArgs() ([]string, error) {
	// For etcd3, we disable the etcd2 endpoint
	// The etcd2 endpoint runs a weird "second copy" of etcd
	if f.Version.IsV2() {
		f.EnableV2 = false
	}

	if f.Quarantined {
		f.ListenClientURLs.Port = QuarantinedClientPort
		f.AdvertiseClientURLs.Port = QuarantinedClientPort
	} else {
		f.ListenClientURLs.Port = ClientPort
		f.AdvertiseClientURLs.Port = ClientPort
	}

	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return meta.BuildArgumentListFromMap(m, nil), nil
}

func (p *EtcdFlags) NewClient() (etcdclient.EtcdClient, error) {
	clientUrls := []string{""}
	if p.Quarantined {
		clientUrls = nil
	}
	return etcdclient.NewClient(string(p.Version), clientUrls)
}
