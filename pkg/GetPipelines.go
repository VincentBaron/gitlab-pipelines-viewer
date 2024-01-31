package pkg

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/xanzy/go-gitlab"
)

type GetPipelinesParams struct {
	Me         bool
	ProjectIDs []string
}

type PipelineProject struct {
	Pipeline    *gitlab.PipelineInfo
	ProjectName string
}

func GetPipelines(_ *gitlab.Client, params GetPipelinesParams) {
	err := godotenv.Load("/Users/vincentbaron/personal/ceyes/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")

	client, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.side.co"))
	if err != nil {
		log.Fatal(err)
	}

	var allPipelines []PipelineProject

	for _, id := range params.ProjectIDs {
		var options gitlab.ListProjectPipelinesOptions
		if params.Me {
			user, _, err := client.Users.CurrentUser()
			if err != nil {
				log.Fatal(err)
			}
			options = gitlab.ListProjectPipelinesOptions{Username: gitlab.String(user.Username)}
		}
		pipelines, _, err := client.Pipelines.ListProjectPipelines(id, &options)
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
		if len(commitMessage) > 25 {
			commitMessage = commitMessage[:22] + "..."
		}
		color.New(color.FgHiBlack).Printf("|📝 %-25s", commitMessage)

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
