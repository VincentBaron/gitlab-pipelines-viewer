package models

import (
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type GetPipelinesParams struct {
	Me         bool
	ProjectIDs []string
}

type Config struct {
	Projects []int `yaml:"projects"`
}

type App struct {
	Client  *gitlab.Client
	RootCmd *cobra.Command
}
