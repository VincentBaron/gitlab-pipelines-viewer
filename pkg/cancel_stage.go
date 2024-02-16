package pkg

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xanzy/go-gitlab"
)

func CancelStage(client *gitlab.Client, args []string) {
	pipelineID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
		return
	}
	stageName := args[1]

	jobs, _, err := client.Jobs.ListPipelineJobs(client, pipelineID, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, job := range jobs {
		if job.Stage == stageName {
			_, _, err := client.Jobs.CancelJob(client, job.ID, nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Cancelled job %d in stage %s\n", job.ID, stageName)
		}
	}
}
