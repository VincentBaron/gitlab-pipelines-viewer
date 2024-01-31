package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"io/ioutil"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gitlab "github.com/xanzy/go-gitlab"

	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "ci",
	Short: "CI/CD pipeline viewer",
}

var removeProject = &cobra.Command{
	Use:   "remove",
	Short: "remove a project from the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeProjectFromConfig(args[0])
	},
}

type Config struct {
	Projects []int `yaml:"projects"`
}

func removeProjectFromConfig(projectName string) {
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

	// Remove the project ID from the config
	for i, id := range config.Projects {
		if id == projects[0].ID {
			config.Projects = append(config.Projects[:i], config.Projects[i+1:]...)
			break
		}
	}

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

	rootCmd.AddCommand(cmd.allCmd, projectCmd, addProject, removeProject, meCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
