package options

import (
	"github.com/etcd-manager/etcd-discovery/pkg/config"
	"github.com/etcd-manager/etcd-discovery/pkg/controller"
	"github.com/spf13/pflag"
)

type EtcdOptions struct {
	ClusterName     string
	ClusterSize     int
	BackupStorePath string

	InitialClusterState config.ClusterState
	InitialCluster      map[string]string
}

func NewEtcdOptions() *EtcdOptions {
	opts := &EtcdOptions{
		InitialClusterState: config.ClusterStateNew,
	}
	return opts
}

func (s *EtcdOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ClusterName, "etcd-cluster-name", s.ClusterName, "Name of cluster")
	fs.IntVar(&s.ClusterSize, "etcd-cluster-size", s.ClusterSize, "Size of cluster size")

	fs.StringVar(&s.BackupStorePath, "backup-store", s.BackupStorePath, "Backup store location")

	fs.StringToStringVar(&s.InitialCluster, "initial-cluster", s.InitialCluster, "Initial cluster configuration")
	fs.Var(&s.InitialClusterState, "initial-cluster-state", "Initial cluster state")
}

func (s *EtcdOptions) Validate() []error {
	return nil
}

func (s *EtcdOptions) ApplyTo(cfg *controller.EtcdConfig) error {
	var err error

	return err
}
