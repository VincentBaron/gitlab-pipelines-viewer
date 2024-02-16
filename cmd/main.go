package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vincentbaron/ceyes/models"
	"github.com/vincentbaron/ceyes/pkg"
	"github.com/xanzy/go-gitlab"
)

var rootCmd = &cobra.Command{
	Use:   "ci",
	Short: "CI/CD pipeline viewer",
}

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
		pkg.GetPipelines(client, models.GetPipelinesParams{ProjectIDs: projects})
	},
}

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
		pkg.GetPipelines(client, models.GetPipelinesParams{
			ProjectIDs: projects,
			Me:         true,
		})
	},
}

var projectCmd = &cobra.Command{
	Use:   "project [name]",
	Short: "Display the pipeline for the specified project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		// token := viper.GetString("token")

		pkg.GetPipelines(nil, models.GetPipelinesParams{ProjectIDs: []string{projectID}})
	},
}

var addProject = &cobra.Command{
	Use:   "add",
	Short: "add a project to the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.AddProjectToConfig(args[0])
	},
}

var removeProject = &cobra.Command{
	Use:   "remove",
	Short: "remove a project from the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.RemoveProjectFromConfig(args[0])
	},
}

func main() {
	// ASCII Art
	data, err := ioutil.ReadFile("/Users/vincentbaron/personal/ceyes/banner.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	lines := strings.Split(string(data), "\n")

	// Print lines 3 to 10
	fmt.Println("\n\n")
	for i := 2; i < 10; i++ {
		fmt.Println(lines[i])
	}
	fmt.Println("\n\n")

	// Rest of your code
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/Users/vincentbaron/personal/ceyes")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	rootCmd.AddCommand(allCmd, projectCmd, addProject, removeProject, meCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
