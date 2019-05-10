package function

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type outputType struct {
	Total  int            `json:"total"`
	Images map[string]int `json:"images"`
}

type dockerHubOrgStatsResult struct {
	User      string `json:"user"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	StarCount int    `json:"star_count"`
	PullCount int    `json:"pull_count"`
}

func consolidate(r []dockerHubOrgStatsResult) outputType {

	var consolidatedRes = make(map[string]int)
	var total int

	for _, image := range r {
		total += image.PullCount
		if _, exists := consolidatedRes[image.Name]; exists {
			consolidatedRes[image.Name] += image.PullCount
		} else {
			consolidatedRes[image.Name] = image.PullCount
		}
	}
	return outputType{Total: total,
		Images: consolidatedRes}
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
	consolidatedResult := consolidate(stats)

	if _, exists := consolidatedResult.Images[image]; exists {
		retStats, marshErr = json.Marshal(consolidatedResult.Images[image])
	} else {
		retStats, marshErr = json.Marshal(consolidatedResult)
	}

	if marshErr != nil {
		log.Fatalln("unable to marshal results slice")
	}

	return fmt.Sprintf("%s", string(retStats))
}
