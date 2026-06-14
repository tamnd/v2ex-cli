package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/v2ex-cli/v2ex"
)

func (a *App) repliesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "replies <topic_id>",
		Short: "List replies for a V2EX topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID := args[0]
			a.progressf("fetching replies for topic %s...", topicID)
			replies, err := a.client.Replies(cmd.Context(), topicID)
			if err != nil {
				return mapFetchErr(err)
			}
			rows := make([]v2ex.ReplyRow, 0, len(replies))
			for i, r := range replies {
				rows = append(rows, v2ex.ReplyToRow(r, i+1))
			}
			return a.renderOrEmpty(rows, len(rows))
		},
	}
}
