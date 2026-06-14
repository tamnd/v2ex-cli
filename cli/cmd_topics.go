package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/v2ex-cli/v2ex"
)

// topicsCmd builds the hot/latest commands from a shared template.
func (a *App) topicsCmd(name, short, endpoint string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			a.progressf("fetching %s topics...", name)
			topics, err := a.client.Topics(cmd.Context(), endpoint)
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
