package controller

import (
	"net"

	api "github.com/etcd-manager/etcd-discovery/apis/discovery/v1alpha1"
	"github.com/etcd-manager/etcd-discovery/pkg/config"
)

type EtcdConfig struct {
	ID               api.PeerID
	AdvertiseAddress net.IP

	ClusterName     string
	ClusterSize     int
	BackupStorePath string
	DataDir         string

	InitialClusterState config.ClusterState
	InitialCluster      map[string]string
}

func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{}
}

func (c *EtcdConfig) New() (*EtcdController, error) {
	return &EtcdController{}, nil
}
