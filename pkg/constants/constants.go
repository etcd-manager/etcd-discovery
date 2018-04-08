package constants

import (
	"fmt"

	"k8s.io/kubernetes/pkg/util/version"
)

const (
	// Etcd defines variable used internally when referring to etcd component
	Etcd = "etcd"

	// EtcdCACertAndKeyBaseName defines etcd's CA certificate and key base name
	EtcdCACertAndKeyBaseName = "etcd/ca"
	// EtcdCACertName defines etcd's CA certificate name
	EtcdCACertName = "etcd/ca.crt"
	// EtcdCAKeyName defines etcd's CA key name
	EtcdCAKeyName = "etcd/ca.key"

	// EtcdServerCertAndKeyBaseName defines etcd's server certificate and key base name
	EtcdServerCertAndKeyBaseName = "etcd/server"
	// EtcdServerCertName defines etcd's server certificate name
	EtcdServerCertName = "etcd/server.crt"
	// EtcdServerKeyName defines etcd's server key name
	EtcdServerKeyName = "etcd/server.key"
	// EtcdServerCertCommonName defines etcd's server certificate common name (CN)
	EtcdServerCertCommonName = "kube-etcd"

	// EtcdPeerCertAndKeyBaseName defines etcd's peer certificate and key base name
	EtcdPeerCertAndKeyBaseName = "etcd/peer"
	// EtcdPeerCertName defines etcd's peer certificate name
	EtcdPeerCertName = "etcd/peer.crt"
	// EtcdPeerKeyName defines etcd's peer key name
	EtcdPeerKeyName = "etcd/peer.key"
	// EtcdPeerCertCommonName defines etcd's peer certificate common name (CN)
	EtcdPeerCertCommonName = "kube-etcd-peer"

	// EtcdHealthcheckClientCertAndKeyBaseName defines etcd's healthcheck client certificate and key base name
	EtcdHealthcheckClientCertAndKeyBaseName = "etcd/healthcheck-client"
	// EtcdHealthcheckClientCertName defines etcd's healthcheck client certificate name
	EtcdHealthcheckClientCertName = "etcd/healthcheck-client.crt"
	// EtcdHealthcheckClientKeyName defines etcd's healthcheck client key name
	EtcdHealthcheckClientKeyName = "etcd/healthcheck-client.key"
	// EtcdHealthcheckClientCertCommonName defines etcd's healthcheck client certificate common name (CN)
	EtcdHealthcheckClientCertCommonName = "kube-etcd-healthcheck-client"

	// kubeControllerManagerAddressArg represents the address argument of the kube-controller-manager configuration.
	kubeControllerManagerAddressArg = "address"

	// kubeSchedulerAddressArg represents the address argument of the kube-scheduler configuration.
	kubeSchedulerAddressArg = "address"

	// etcdListenClientURLsArg represents the listen-client-urls argument of the etcd configuration.
	EtcdListenClientURLsArg = "listen-client-urls"

	// DefaultEtcdVersion indicates the default etcd version that kubeadm uses
	DefaultEtcdVersion = "3.1.12"
)

var (
	// SupportedEtcdVersion lists officially supported etcd versions with corresponding kubernetes releases
	SupportedEtcdVersion = map[uint8]string{
		9:  "3.1.12",
		10: "3.1.12",
		11: "3.1.12",
	}
)

// EtcdSupportedVersion returns officially supported version of etcd for a specific kubernetes release
// if passed version is not listed, the function returns nil and an error
func EtcdSupportedVersion(versionString string) (*version.Version, error) {
	kubernetesVersion, err := version.ParseSemantic(versionString)
	if err != nil {
		return nil, err
	}

	if etcdStringVersion, ok := SupportedEtcdVersion[uint8(kubernetesVersion.Minor())]; ok {
		etcdVersion, err := version.ParseSemantic(etcdStringVersion)
		if err != nil {
			return nil, err
		}
		return etcdVersion, nil
	}
	return nil, fmt.Errorf("Unsupported or unknown kubernetes version")
}
