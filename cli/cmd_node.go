package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/v2ex-cli/v2ex"
)

func (a *App) nodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "node <name>",
		Short: "Show a V2EX node by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			a.progressf("fetching node %s...", name)
			n, err := a.client.Node(cmd.Context(), name)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.render(v2ex.NodeToDetails(*n))
		},
	}
}
