package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tangx/gitlab-mergebot/cmd/gitlab-mergebot/global"
	"github.com/tangx/gitlab-mergebot/pkg/launcher"
	"github.com/tangx/gitlab-mergebot/version"
)

var rootCmd = &cobra.Command{
	Use:  "gitlab-mergebot",
	Long: version.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	Run: func(cmd *cobra.Command, args []string) {
		launch()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func launch() {
	launchor := launcher.New(10)
	launchor.WithFuncs(global.MergebotRun)
	launchor.Launch()
}
