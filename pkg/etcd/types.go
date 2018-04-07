package etcd

import "os"

type Process interface {
	Type() ProcessType
	Start() error
	Stop() error
	ExitState() (error, *os.ProcessState)
}
