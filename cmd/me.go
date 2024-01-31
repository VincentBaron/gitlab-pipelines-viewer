package cmd

import (
	"log"

	"github.com/VincentBaron/ceyes/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Display my ongoing pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		projects := viper.GetStringSlice("projects")
		token := viper.GetString("token")

		client, err := gitlab.NewClient(token)
		if err != nil {
			log.Fatal(err)
		}
		pkg.GetPipelines(client, pkg.GetPipelinesParams{
			ProjectIDs: projects,
			Me:         true,
		})
	},
}
