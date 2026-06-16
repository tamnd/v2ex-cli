package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/v2ex-cli/v2ex"
)

// nodeTopicsCmd returns topics for a given node by name.
func (a *App) nodeTopicsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "topics <node_name>",
		Short: "List recent topics in a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			a.progressf("fetching topics for node %s...", name)
			topics, err := a.client.TopicsByNode(cmd.Context(), name)
			if err != nil {
				return mapFetchErr(err)
			}
			n := a.effectiveLimit(len(topics))
			if n > len(topics) {
				n = len(topics)
			}
			rows := make([]v2ex.TopicRow, 0, n)
			for i, t := range topics[:n] {
				rows = append(rows, v2ex.TopicToRow(t, i+1))
			}
			return a.renderOrEmpty(rows, len(rows))
		},
	}
}
