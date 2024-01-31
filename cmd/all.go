package cmd

import (
	"log"

	"github.com/VincentBaron/ceyes/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Display all ongoing pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		projects := viper.GetStringSlice("projects")
		token := viper.GetString("token")

		client, err := gitlab.NewClient(token)
		if err != nil {
			log.Fatal(err)
		}
		pkg.GetPipelines(client, pkg.GetPipelinesParams{ProjectIDs: projects})
	},
}
