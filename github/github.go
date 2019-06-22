package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Project struct {
	Name            string `xml:"name,attr"`
	Activity        string `xml:"activity,attr"`
	LastBuildTime   string `xml:"lastBuildTime,attr"`
	URL             string `xml:"webUrl,attr"`
	LastBuildStatus string `xml:"lastBuildStatus,attr"`
}

func Activity(state string) string {
	activity := "Sleeping"
	switch state {
	case "failure":
		activity = "Failure"
	case "pending":
		activity = "Building"
	case "success":
		activity = "Success"
	}
	return activity
}

// GetStatus - Query Github API for commit status of master
func GetProject(repo string, token string) Project {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.github.com/repos/%s/commits/master/status", repo)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Add("User-Agent", "cchub")
	resp, err := client.Do(req)

	if err != nil {
		println(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	statuses := result["statuses"].([]interface{})
	lastStatus := statuses[0].(map[string]interface{})
	lastBuildTime, _ := time.Parse(time.RFC3339, lastStatus["updated_at"].(string))
	repoURL := fmt.Sprintf("https://github.com/%s/commits/master", repo)
	lastBuildState := lastStatus["state"].(string)

	return Project{
		Name:            repo,
		Activity:        Activity(result["state"].(string)),
		LastBuildTime:   lastBuildTime.Format(time.RFC3339),
		URL:             repoURL,
		LastBuildStatus: Activity(lastBuildState),
	}
}
