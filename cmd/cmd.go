package cmd

import "github.com/spf13/cobra"

var usage = `skool-mn
	- server: start server
@skool-mn
`

// New rootCommand
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "server",
		Short:        "skool-mn",
		Long:         usage,
		SilenceUsage: true,
	}
	cmd.AddCommand(newServerCmd())
	return cmd
}
