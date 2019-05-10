package function

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type dockerHubOrgStatsResult struct {
	User      string `json:"user"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	StarCount int    `json:"star_count"`
	PullCount int    `json:"pull_count"`
}

func consolidate(r []dockerHubOrgStatsResult) map[string]int {

	var consolidatedRes = make(map[string]int)

	for _, image := range r {
		if _, exists := consolidatedRes[image.Name]; exists {
			consolidatedRes[image.Name] += image.PullCount
		} else {
			consolidatedRes[image.Name] = image.PullCount
		}
	}
	return consolidatedRes
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {

	var envVal string
	if value, exists := os.LookupEnv(name); exists {
		envVal = value
	}

	if len(envVal) == 0 {
		return defaultVal
	}
	return strings.Split(envVal, sep)
}

func getStats(orgs []string) []dockerHubOrgStatsResult {

	var results []dockerHubOrgStatsResult
	var statsResponse dockerHubOrgStats

	for _, org := range orgs {

		pageNo := 1
		for {
			statsResponse = requestStats(org, pageNo)
			results = append(results, statsResponse.Results...)
			if statsResponse.isLast() {
				break
			}
			pageNo++
		}
	}
	return results
}

// Handle a serverless request
func Handle(req []byte) string {

	var marshErr error
	var retStats []byte
	image := string(req)

	orgs := getEnvAsSlice("orgs", []string{"rgee0"}, ",")
	stats := getStats(orgs)
	consolidatedResults := consolidate(stats)

	if _, exists := consolidatedResults[image]; exists {
		retStats, marshErr = json.Marshal(consolidatedResults[image])
	} else {
		retStats, marshErr = json.Marshal(consolidatedResults)
	}

	if marshErr != nil {
		log.Fatalln("unable to marshal results slice")
	}

	return fmt.Sprintf("%s", string(retStats))
}
