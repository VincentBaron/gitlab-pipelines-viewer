package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gitlab "github.com/xanzy/go-gitlab"
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

		for _, projectID := range projects {
			getPipelines(client, projectID)
		}
	},
}

var projectCmd = &cobra.Command{
	Use:   "project [name]",
	Short: "Display the pipeline for the specified project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		// token := viper.GetString("token")

		getPipelines(nil, projectID)
	},
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	rootCmd.AddCommand(allCmd, projectCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getPipelines(_ *gitlab.Client, projectID string) {
	client, err := gitlab.NewClient("glpat-mWbmz1Mb6B2q5xrqvpnH", gitlab.WithBaseURL("https://gitlab.side.co"))
	if err != nil {
		log.Fatal(err)
	}

	pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, &gitlab.ListProjectPipelinesOptions{})
	if err != nil {
		fmt.Println("Oops! Something went wrong 🙈")
		log.Fatal(err)
	}

	project, _, err := client.Projects.GetProject(projectID, &gitlab.GetProjectOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Define the order of stages
	stageOrder := []string{"check", "sonar", "build", "scan", "squad1", "squad2", "dev", "prod", "roll", "secrets"}

	for _, pipeline := range pipelines {
		// Get jobs of the pipeline
		jobs, _, err := client.Jobs.ListPipelineJobs(projectID, pipeline.ID, &gitlab.ListJobsOptions{})
		if err != nil {
			log.Fatal(err)
		}

		// Group jobs by stage
		stageJobs := make(map[string][]*gitlab.Job)
		hasFailedOrRunningStage := false
		for _, job := range jobs {
			stageJobs[job.Stage] = append(stageJobs[job.Stage], job)
			if job.Status == "failed" || job.Status == "running" || job.Status == "pending" {
				hasFailedOrRunningStage = true
			}
		}

		// Print pipeline and stage information if there's at least one failed or running stage
		if hasFailedOrRunningStage {
			var colorPrinter func(format string, a ...interface{})
			pipelineStatus := "other"
			for _, jobs := range stageJobs {
				for _, job := range jobs {
					if job.Status == "failed" {
						pipelineStatus = "failed"
						break
					} else if job.Status == "running" || job.Status == "pending" {
						pipelineStatus = "running"
					} else if job.Status == "success" {
						pipelineStatus = "success"
					}
				}
				if pipelineStatus == "failed" {
					break
				}
			}

			switch pipelineStatus {
			case "success":
				colorPrinter = color.New(color.FgGreen).PrintfFunc()
			case "failed":
				colorPrinter = color.New(color.FgRed).PrintfFunc()
			default:
				colorPrinter = color.New(color.FgYellow).PrintfFunc()
			}

			colorPrinter("🚀 %s | 🪵 %s | ", project.Name, pipeline.Ref)

			// Iterate over stages in the defined order
			for _, stage := range stageOrder {
				jobs, ok := stageJobs[stage]
				if !ok {
					continue
				}

				stageStatus := "other"
				for _, job := range jobs {
					if job.Status == "failed" {
						stageStatus = "failed"
						break
					} else if job.Status == "running" || job.Status == "pending" {
						stageStatus = "running"
					} else if job.Status == "success" {
						stageStatus = "success"
					}
				}
				switch stageStatus {
				case "running":
					fmt.Printf("🔧 %s ➡️ ", stage)
				case "failed":
					fmt.Printf("💥 %s ➡️ ", stage)
				case "success":
					fmt.Printf("✅ %s ➡️ ", stage)
				default:
					fmt.Printf("⏳ %s ➡️ ", stage)
				}
			}
			fmt.Println()
		}
	}
}
