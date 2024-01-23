package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
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

        for _, projectID := range projects {
            getPipelines(token, projectID)
        }
    },
}

var projectCmd = &cobra.Command{
    Use:   "project [name]",
    Short: "Display the pipeline for the specified project",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        projectID := args[0]
        token := viper.GetString("token")

        getPipelines(token, projectID)
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

func getPipelines(token string, projectID string) {
    url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/pipelines", projectID)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("PRIVATE-TOKEN", token)

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    fmt.Printf("Pipelines for project %s:\n%s\n", projectID, string(body))
}