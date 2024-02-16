package pkg

import (
	"io/ioutil"
	"log"

	"github.com/vincentbaron/ceyes/models"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

func AddProjectToConfig(client *gitlab.Client, projectName string) {
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

	var config models.Config

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
