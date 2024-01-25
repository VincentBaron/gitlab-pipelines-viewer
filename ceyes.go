package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"io/ioutil"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gitlab "github.com/xanzy/go-gitlab"

	"gopkg.in/yaml.v2"
)

type PipelineProject struct {
	Pipeline    *gitlab.PipelineInfo
	ProjectName string
}

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
		getPipelines(client, projects)
	},
}

var projectCmd = &cobra.Command{
	Use:   "project [name]",
	Short: "Display the pipeline for the specified project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		// token := viper.GetString("token")

		getPipelines(nil, []string{projectID})
	},
}

var addProject = &cobra.Command{
	Use:   "add",
	Short: "add a project to the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addProjectToConfig(args[0])
	},
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

func main() {
	// ASCII Art
	data, err := ioutil.ReadFile("banner.txt")
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

	rootCmd.AddCommand(allCmd, projectCmd, addProject, removeProject)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getPipelines(_ *gitlab.Client, projectIDs []string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")

	client, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.side.co"))
	if err != nil {
		log.Fatal(err)
	}

	var allPipelines []PipelineProject

	for _, id := range projectIDs {
		pipelines, _, err := client.Pipelines.ListProjectPipelines(id, &gitlab.ListProjectPipelinesOptions{})
		if err != nil {
			log.Fatal(err)
		}

		project, _, err := client.Projects.GetProject(id, &gitlab.GetProjectOptions{})
		if err != nil {
			log.Fatal(err)
		}

		for _, pipeline := range pipelines {
			intID, err := strconv.Atoi(id)
			if err != nil {
				log.Fatal(err)
			}
			pipeline.ProjectID = intID

			allPipelines = append(allPipelines, PipelineProject{Pipeline: pipeline, ProjectName: project.Name})
		}
	}

	// Define the order of stages
	stageOrder := []string{"check", "sonar", "build", "scan", "squad1", "squad2", "dev", "prod", "roll", "secrets"}

	// Get the current time in Paris
	location, _ := time.LoadLocation("Europe/Paris")
	now := time.Now().In(location)

	// Set the start of the day to 8 AM
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, location)

	// Sort pipelines by creation date in descending order
	sort.SliceStable(allPipelines, func(i, j int) bool {
		return allPipelines[i].Pipeline.UpdatedAt.After(*allPipelines[j].Pipeline.UpdatedAt)
	})

	// Calculate maximum lengths
	maxProjectNameLen := 0
	maxRefLen := 0
	for _, pipeline := range allPipelines {
		if pipeline.Pipeline.UpdatedAt.Before(startOfDay) {
			continue
		}
		if len(pipeline.ProjectName) > maxProjectNameLen {
			maxProjectNameLen = len(pipeline.ProjectName)
		}
		if len(pipeline.Pipeline.Ref) > maxRefLen {
			maxRefLen = len(pipeline.Pipeline.Ref)
		}
	}

	// Print data
	for _, pipeline := range allPipelines {
		// Only consider pipelines created after 8 AM
		if pipeline.Pipeline.UpdatedAt.Before(startOfDay) {
			continue
		}

		// Get jobs of the pipeline
		jobs, _, err := client.Jobs.ListPipelineJobs(strconv.Itoa(pipeline.Pipeline.ProjectID), pipeline.Pipeline.ID, &gitlab.ListJobsOptions{})
		if err != nil {
			log.Fatal(err)
		}

		// Group jobs by stage
		stageJobs := make(map[string][]*gitlab.Job)
		for _, job := range jobs {
			stageJobs[job.Stage] = append(stageJobs[job.Stage], job)
		}

		var colorPrinter func(format string, a ...interface{})

		pipelineStatus := pipeline.Pipeline.Status

		switch pipelineStatus {
		case "success":
			colorPrinter = color.New(color.FgGreen).PrintfFunc()
		case "failed":
			colorPrinter = color.New(color.FgRed).PrintfFunc()
		case "running", "pending":
			colorPrinter = color.New(color.FgYellow).PrintfFunc()
		default:
			colorPrinter = color.New(color.FgWhite).PrintfFunc()
		}

		colorPrinter("🚀 %-*s", maxProjectNameLen, pipeline.ProjectName)
		color.New(color.FgBlue).Printf("|🪵  %-*s", maxRefLen, pipeline.Pipeline.Ref)

		// Get the commit message
		commit, _, err := client.Commits.GetCommit(strconv.Itoa(pipeline.Pipeline.ProjectID), pipeline.Pipeline.SHA)
		if err != nil {
			log.Fatal(err)
		}
		commitMessage := commit.Title
		if len(commitMessage) > 30 {
			commitMessage = commitMessage[:27] + "..."
		}
		color.New(color.FgHiBlack).Printf("|📝 %-30s", commitMessage)

		// Iterate over stages in the defined order
		for _, stage := range stageOrder {
			jobs, ok := stageJobs[stage]
			if !ok {
				continue
			}

			stageStatus := "other"
			allJobsManual := true
			allJobsPassed := true
			for _, job := range jobs {
				if job.Status != "manual" {
					allJobsManual = false
				}
				if job.Status != "success" {
					allJobsPassed = false
				}
				if job.Status == "running" {
					stageStatus = "running"
					break
				} else if job.Status == "pending" {
					stageStatus = "pending"
					break
				} else if job.Status == "failed" {
					if job.AllowFailure {
						stageStatus = "allowed to fail"
					} else {
						stageStatus = "failed"
					}
					break
				}
			}
			if allJobsManual {
				stageStatus = "manual"
			} else if allJobsPassed {
				stageStatus = "passed"
			}

			// fmt.Println(stageStatus)
			switch stageStatus {
			case "running":
				color.New(color.FgMagenta).Printf("|🟣 %-6s", stage)
			case "failed":
				color.New(color.FgRed).Printf("|💥 %-6s", stage)
			case "passed":
				color.New(color.FgGreen).Printf("|✅ %-6s", stage)
			case "other":
				color.New(color.FgBlue).Printf("|👋 %-6s", stage)
			case "pending":
				color.New(color.FgWhite).Printf("|⏳ %-6s", stage)
			case "allowed to fail":
				color.New(color.FgYellow).Printf("|🟨 %-6s", stage)
			default:
				color.New(color.FgHiBlack).Printf("|❓ %-6s", stage)
			}
		}
		fmt.Println("\n-------------------------------------------------------------------------------------------------------------------------------------------------------")
	}
}
