package cmds

import (
	"fmt"

	"github.com/appscode/go/log"
	"github.com/appscode/go/net"
	"github.com/appscode/kutil/tools/certstore"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/cert"
)

const organization = "system:etcd"

func NewCmdConfigure() *cobra.Command {
	var (
		certDir = "etcd.local.config/certificates"
		addr    = "127.0.0.1"
	)
	cmd := &cobra.Command{
		Use:               "configure",
		Short:             "Configure certs for etcd-discovery",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			externalIPs, internalIPs, err := net.HostIPs()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("extenal-ips:", externalIPs)
			fmt.Println("internal-ips:", internalIPs)

			err = prepareServerCerts(certDir, "client", addr)
			if err != nil {
				log.Fatal(err)
			}
			err = preparePeerCerts(certDir, "peer", addr)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&certDir, "cert-dir", certDir, "Path to directory where pki files are stored.")
	cmd.Flags().StringVar(&addr, "addr", addr, "Address of server ip")
	return cmd
}

func prepareServerCerts(certDir, ca, addr string) error {
	store, err := certstore.NewCertStore(afero.NewOsFs(), certDir, organization)
	if err != nil {
		return errors.Wrap(err, "failed to create certificate store.")
	}
	err = store.InitCA(ca)
	if err != nil {
		return errors.Wrap(err, "failed to init ca.")
	}
	crt, key, err := store.NewServerCertPair(addr, cert.AltNames{})
	if err != nil {
		return err
	}
	return store.WriteBytes(addr, crt, key)
}

func preparePeerCerts(certDir, ca, addr string) error {
	store, err := certstore.NewCertStore(afero.NewOsFs(), certDir, organization)
	if err != nil {
		return errors.Wrap(err, "failed to create certificate store.")
	}
	err = store.InitCA(ca)
	if err != nil {
		return errors.Wrap(err, "failed to init ca.")
	}
	crt, key, err := store.NewPeerCertPair(addr, cert.AltNames{})
	if err != nil {
		return err
	}
	return store.WriteBytes(addr, crt, key)
}
