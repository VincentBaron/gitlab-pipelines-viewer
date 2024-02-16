package pkg

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/vincentbaron/ceyes/models"
	"github.com/xanzy/go-gitlab"
)

type PipelineProject struct {
	Pipeline    *gitlab.PipelineInfo
	ProjectName string
}

func GetPipelines(client *gitlab.Client, params models.GetPipelinesParams) {
	var allPipelines []PipelineProject

	// Get the current time in Paris
	location, _ := time.LoadLocation("Europe/London")
	now := time.Now().In(location)

	// Set the start of the day to 8 AM
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, location)

	for _, id := range params.ProjectIDs {
		options := gitlab.ListProjectPipelinesOptions{
			UpdatedAfter: gitlab.Ptr(startOfDay),
		}
		if params.Me {
			user, _, err := client.Users.CurrentUser()
			if err != nil {
				log.Fatal(err)
			}
			options.Username = gitlab.Ptr(user.Username)
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

	// Sort pipelines by creation date in descending order
	sort.SliceStable(allPipelines, func(i, j int) bool {
		return allPipelines[i].Pipeline.UpdatedAt.After(*allPipelines[j].Pipeline.UpdatedAt)
	})

	// Calculate maximum lengths
	maxProjectNameLen := 0
	maxRefLen := 0
	for _, pipeline := range allPipelines {
		if len(pipeline.ProjectName) > maxProjectNameLen {
			maxProjectNameLen = len(pipeline.ProjectName)
		}
		if len(pipeline.Pipeline.Ref) > maxRefLen {
			maxRefLen = len(pipeline.Pipeline.Ref)
		}

	}

	// Print data
	for _, pipeline := range allPipelines {
		// Get jobs of the pipeline
		jobs, _, err := client.Jobs.ListPipelineJobs(strconv.Itoa(pipeline.Pipeline.ProjectID), pipeline.Pipeline.ID, &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100},
		})
		if err != nil {
			log.Fatal(err)
		}

		// Define the order of jobs
		jobOrder := map[string]int{
			"check":               1,
			"sonarqube-check":     2,
			"build":               3,
			"vulnerability_check": 4,
			"deploy to squad-1":   5,
			"deploy to squad-2":   6,
			"deploy to dev":       7,
			"deploy to prod":      8,
			"sentry_release":      9,
		}

		// Filter jobs
		filteredJobs := make([]*gitlab.Job, 0)
		for _, job := range jobs {
			if _, ok := jobOrder[job.Name]; ok {
				filteredJobs = append(filteredJobs, job)
			}
		}

		// Sort jobs by the defined order
		sort.SliceStable(filteredJobs, func(i, j int) bool {
			return jobOrder[filteredJobs[i].Name] < jobOrder[filteredJobs[j].Name]
		})

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

		// Print pipeline and commit info
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

		// Iterate over jobs and print their status
		for _, job := range filteredJobs {
			stageName := job.Stage
			if stageName == "production" {
				stageName = "prod"
			}
			switch job.Status {
			case "running":
				color.New(color.FgMagenta).Printf("|🟣 %-7s", stageName)
			case "failed":
				color.New(color.FgRed).Printf("|💥 %-7s", stageName)
			case "success":
				color.New(color.FgGreen).Printf("|✅ %-7s", stageName)
			case "other":
				color.New(color.FgBlue).Printf("|👋 %-7s", stageName)
			case "pending":
				color.New(color.FgWhite).Printf("|⏳ %-7s", stageName)
			case "canceled":
				color.New(color.FgYellow).Printf("|⏹️  %-7s", stageName)
			case "skipped":
				color.New(color.FgYellow).Printf("|⏭️  %-7s", stageName)
			default:
				color.New(color.FgHiBlack).Printf("|❓ %-7s", stageName)
			}
		}
		fmt.Println("\n-------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	}
}
