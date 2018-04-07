package etcd

import (
	"bytes"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	ioutils "github.com/appscode/go/ioutil"
	api "github.com/etcd-manager/etcd-discovery/apis/discovery/v1alpha1"
	"github.com/golang/glog"
)

// PersistentPeerId reads the id from the base directory, creating and saving it if it does not exists
func PersistentPeerId(basedir string) (api.PeerId, error) {
	idFile := filepath.Join(basedir, "myid")

	b, err := ioutil.ReadFile(idFile)
	if err != nil {
		if os.IsNotExist(err) {
			token := randomToken()
			glog.Infof("Self-assigned new identity: %q", token)
			if err := ioutils.WriteFile(idFile, bytes.NewBufferString(token), 0644); err != nil {
				return "", fmt.Errorf("error creating id file %q: %v", idFile, err)
			}
		} else {
			return "", fmt.Errorf("error reading id file %q: %v", idFile, err)
		}
	}

	uniqueID := api.PeerId(string(b))
	return uniqueID, nil
}

func randomToken() string {
	b := make([]byte, 16, 16)
	_, err := io.ReadFull(crypto_rand.Reader, b)
	if err != nil {
		glog.Fatalf("error generating random token: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
