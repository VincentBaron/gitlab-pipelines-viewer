package pkg

import (
	"io/ioutil"
	"log"

	"github.com/vincentbaron/ceyes/models"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

func RemoveProjectFromConfig(client *gitlab.Client, projectName string) {
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
