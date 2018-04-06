package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Run etcd discovery server",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello! Using etcd-discovery")
		},
	}

	return cmd
}
