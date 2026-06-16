package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/v2ex-cli/v2ex"
)

// allNodesCmd lists every V2EX node.
func (a *App) allNodesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-nodes",
		Short: "List all V2EX nodes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			a.progressf("fetching all nodes...")
			nodes, err := a.client.AllNodes(cmd.Context())
			if err != nil {
				return mapFetchErr(err)
			}
			n := a.effectiveLimit(len(nodes))
			if n > len(nodes) {
				n = len(nodes)
			}
			rows := make([]v2ex.NodeRow, 0, n)
			for i, nd := range nodes[:n] {
				rows = append(rows, v2ex.NodeToRow(nd, i+1))
			}
			return a.renderOrEmpty(rows, len(rows))
		},
	}
}
