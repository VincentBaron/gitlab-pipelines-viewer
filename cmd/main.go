package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vincentbaron/ceyes/models"
	"github.com/vincentbaron/ceyes/pkg"
	"github.com/xanzy/go-gitlab"
)

// Global client variable
var client *gitlab.Client

var rootCmd = &cobra.Command{
	Use:   "ci",
	Short: "CI/CD pipeline viewer",
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Display all ongoing pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		projects := viper.GetStringSlice("projects")
		pkg.GetPipelines(client, models.GetPipelinesParams{ProjectIDs: projects})
	},
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Display my ongoing pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		projects := viper.GetStringSlice("projects")
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
		pkg.GetPipelines(client, models.GetPipelinesParams{ProjectIDs: []string{projectID}})
	},
}

var addProject = &cobra.Command{
	Use:   "add",
	Short: "add a project to the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.AddProjectToConfig(client, args[0])
	},
}

var removeProject = &cobra.Command{
	Use:   "remove",
	Short: "remove a project from the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.RemoveProjectFromConfig(client, args[0])
	},
}

var cancelStage = &cobra.Command{
	Use:   "cl [pipeline-id] [stage-name]",
	Short: "Cancel a job by its stage name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.CancelStage(client, args)
	},
}

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// cwd = filepath.Dir(cwd)

	err = godotenv.Load(cwd + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")

	client, err = gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.side.co"))
	if err != nil {
		log.Fatal(err)
	}

	// ASCII Art
	bannerPath := cwd + "/assets/banner.txt"
	data, err := ioutil.ReadFile(bannerPath)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	lines := strings.Split(string(data), "\n")

	// Print lines 3 to 10
	fmt.Println()
	for i := 2; i < 10; i++ {
		fmt.Println(lines[i])
	}
	fmt.Println()

	// Rest of your code
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cwd)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	rootCmd.AddCommand(allCmd, projectCmd, addProject, removeProject, meCmd, cancelStage)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
