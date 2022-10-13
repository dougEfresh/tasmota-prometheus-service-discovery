package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasmota-prometheus-service-discovery",
		Short: "golang-cli project template demo application",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprint(cmd.OutOrStderr(), cmd.UsageString())

			return nil
		},
	}

	cmd.AddCommand(newVersionCmd(version)) // version subcommand
	cmd.AddCommand(newMqttCmd())

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}
