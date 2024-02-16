package models

type GetPipelinesParams struct {
	Me         bool
	ProjectIDs []string
}

type Config struct {
	Projects []int `yaml:"projects"`
}
