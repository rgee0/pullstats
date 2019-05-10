package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type dockerHubOrgStats struct {
	Count    int                       `json:"count"`
	Next     string                    `json:"next"`
	Previous string                    `json:"previous"`
	Results  []dockerHubOrgStatsResult `json:"results"`
}

func (d *dockerHubOrgStats) isLast() bool {
	return len(d.Next) == 0
}

func (d *dockerHubOrgStats) hasNext() bool {
	return len(d.Next) > 0
}

func parseOrgStats(response []byte) (dockerHubOrgStats, error) {
	dockerHubOrgStats := dockerHubOrgStats{}
	err := json.Unmarshal(response, &dockerHubOrgStats)
	return dockerHubOrgStats, err
}

func requestStats(org string, pageNo int) dockerHubOrgStats {

	client := http.Client{}
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/?page=%d&page_size=100", org, pageNo)

	res, err := client.Get(url)
	if err != nil {
		log.Fatalln("Unable to reach Docker Hub endpoint.")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("Unable to parse response from server.")
	}

	orgstats := dockerHubOrgStats{}
	json.Unmarshal(body, &orgstats)
	return orgstats
}
