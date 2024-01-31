package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

var addProject = &cobra.Command{
	Use:   "add",
	Short: "add a project to the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addProjectToConfig(args[0])
	},
}

func addProjectToConfig(projectName string) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get GitLab token
	token := os.Getenv("TOKEN")

	// Create GitLab client
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.side.co"))
	if err != nil {
		log.Fatal(err)
	}

	// Search for the project
	options := &gitlab.ListProjectsOptions{
		Search: gitlab.String(projectName),
	}

	projects, _, err := client.Projects.ListProjects(options)
	if err != nil {
		log.Fatal(err)
	}

	if len(projects) == 0 {
		log.Fatalf("No project found with name: %s", projectName)
	}

	// Load the current config
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Add the project ID to the config
	config.Projects = append(config.Projects, projects[0].ID)

	// Save the updated config
	data, err = yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile("config.yaml", data, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
