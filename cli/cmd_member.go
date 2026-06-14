package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/v2ex-cli/v2ex"
)

func (a *App) memberCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "member <username>",
		Short: "Show a V2EX member profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			a.progressf("fetching member %s...", username)
			m, err := a.client.Member(cmd.Context(), username)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.render(v2ex.MemberToDetails(*m))
		},
	}
}
