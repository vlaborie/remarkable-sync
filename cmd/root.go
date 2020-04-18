package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "remarkable-sync",
		Short: "Sync tool for reMarkable paper tablet",
		Long: `Remarkable-sync is a Go applications for syncing external
services to reMarkable paper table, like Wallabag or Miniflux.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
