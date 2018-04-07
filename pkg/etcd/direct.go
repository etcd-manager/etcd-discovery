package etcd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/etcd-manager/etcd-discovery/pkg/config"
	"github.com/golang/glog"
)

// etcdDirect wraps a running etcd process
type etcdDirect struct {
	BinDir string
	cfg    *config.EtcdFlags

	cmd *exec.Cmd

	mutex     sync.Mutex
	exitError error
	exitState *os.ProcessState
}

var _ Process = &etcdDirect{}

func (p *etcdDirect) Type() ProcessType {
	return ProcessTypeDirect
}

// BindirForEtcdVersion returns the directory in which the etcd binary is located, for the specified version
// It returns an error if the specified version cannot be found
func BindirForEtcdVersion(etcdVersion string, cmd string) (string, error) {
	if !strings.HasPrefix(etcdVersion, "v") {
		etcdVersion = "v" + etcdVersion
	}
	binDir := filepath.Join("/opt", "etcd-"+etcdVersion+"-"+runtime.GOOS+"-"+runtime.GOARCH)
	etcdBinary := filepath.Join(binDir, cmd)
	_, err := os.Stat(etcdBinary)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("unknown etcd version (%s not found at %s)", cmd, etcdBinary)
		} else {
			return "", fmt.Errorf("error checking for %s at %s: %v", cmd, etcdBinary, err)
		}
	}
	return binDir, nil
}

func (p *etcdDirect) Start() error {
	c := exec.Command(path.Join(p.BinDir, "etcd"))

	args, err := p.cfg.ToArgs()
	if err != nil {
		return err
	}
	c.Args = args
	glog.Infof("executing command %s %s", c.Path, c.Args)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Start()
	if err != nil {
		return fmt.Errorf("error starting etcd: %v", err)
	}
	p.cmd = c

	go func() {
		processState, err := p.cmd.Process.Wait()
		if err != nil {
			glog.Warningf("etcd exited with error: %v", err)
		}
		p.mutex.Lock()
		p.exitState = processState
		p.exitError = err
		p.mutex.Unlock()
	}()

	return nil
}

func (p *etcdDirect) Stop() error {
	if p.cmd == nil {
		glog.Warningf("received Stop when process not running")
		return nil
	}
	if err := p.cmd.Process.Kill(); err != nil {
		p.mutex.Lock()
		if p.exitState != nil {
			glog.Infof("Exited etcd: %v", p.exitState)
			return nil
		}
		p.mutex.Unlock()
		return fmt.Errorf("failed to kill process: %v", err)
	}

	for {
		glog.Infof("Waiting for etcd to exit")
		p.mutex.Lock()
		if p.exitState != nil {
			exitState := p.exitState
			p.mutex.Unlock()
			glog.Infof("Exited etcd: %v", exitState)
			return nil
		}
		p.mutex.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (p *etcdDirect) ExitState() (error, *os.ProcessState) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.exitError, p.exitState
}
