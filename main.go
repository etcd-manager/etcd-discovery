package main

import (
	"os"
	"runtime"

	"github.com/etcd-manager/etcd-discovery/pkg/cmds"
	"github.com/golang/glog"
	"k8s.io/apiserver/pkg/util/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if err := cmds.NewRootCmd().Execute(); err != nil {
		glog.Fatalln("Error in etcd-discovery:", err)
	}
	os.Exit(0)
}
