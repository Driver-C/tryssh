package prune

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
	"time"
)

const (
	concurrency = 8
	sshTimeout  = 2 * time.Second
)

func NewPruneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Check if all current caches are available and clear the ones that are not available",
		Long:  "Check if all current caches are available and clear the ones that are not available",
		Run: func(cmd *cobra.Command, args []string) {
			auto, _ := cmd.Flags().GetBool("auto")
			concurrencyOpt, _ := cmd.Flags().GetInt("concurrency")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			configuration := config.LoadConfig()
			controller := control.NewPruneController(configuration, auto, timeout, concurrencyOpt)
			controller.PruneCaches()
		},
	}
	cmd.Flags().BoolP(
		"auto", "a", false, "Automatically perform concurrent cache optimization without"+
			" asking for confirmation to delete")
	cmd.Flags().IntP(
		"concurrency", "c", concurrency, "Number of multiple requests to perform at a time")
	cmd.Flags().DurationP("timeout", "t", sshTimeout,
		"SSH timeout when attempting to log in. It can be \"1s\" or \"1m\" or other duration")
	return cmd
}
