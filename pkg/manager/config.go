package manager

import (
	"net"

	api "github.com/etcd-manager/etcd-discovery/apis/discovery/v1alpha1"
	"github.com/etcd-manager/etcd-discovery/pkg/config"
)

type EtcdConfig struct {
	config.EtcdCluster

	ID               api.PeerID
	AdvertiseAddress net.IP
}

func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{}
}

func (c *EtcdConfig) New() (*EtcdManager, error) {
	return &EtcdManager{}, nil
}
