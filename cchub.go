package main

import (
	"cchub/github"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo"
	"github.com/mitchellh/go-homedir"
)

type Projects struct {
	Projects []github.Project `xml:"Project"`
}

func main() {
	e := echo.New()
	config := LoadConfig()

	if config.GithubToken == "" {
		os.Exit(1)
	}

	e.GET("/", func(c echo.Context) error {
		projects := Projects{}

		for _, repo := range config.Repos {
			projects.Projects = append(projects.Projects, github.GetProject(repo, config.GithubToken))
		}

		return c.XML(http.StatusOK, projects)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Port)))
}

type Config struct {
	Port        int
	GithubToken string
	Repos       []string
}

var DefaultConfig = Config{
	Port:        1323,
	GithubToken: "",
	Repos:       []string{},
}

func LoadConfig() *Config {
	homeDir, _ := homedir.Dir()
	configFile := filepath.Join(homeDir, ".cchub.toml")
	conf := DefaultConfig

	if _, err := os.Stat(configFile); err != nil {
		log.Println("config file does not exist")
		return &conf
	}

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		return &conf
	}

	return &conf
}
