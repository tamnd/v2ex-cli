package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/v2ex-cli/v2ex"
)

func (a *App) topicCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "topic <id>",
		Short: "Show a single V2EX topic by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			a.progressf("fetching topic %s...", id)
			t, err := a.client.Topic(cmd.Context(), id)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.render(v2ex.TopicToDetails(*t))
		},
	}
}
