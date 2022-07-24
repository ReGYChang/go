package main

import (
	"encoding/json"
	"io/ioutil"
	"nexdata/pkg/config"
	"nexdata/pkg/database/elasticsearch"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/rs/zerolog/log"
)

func main() {
	//jsonFile, err := os.Open("sensormap_sgroup_version.json")

	jsonFile, err := ioutil.ReadFile("pkg/database/elasticsearch/cmd/sensormap_sgroup_version.json")
	if err != nil {
		log.Fatal().Msgf("Error when opening file: ", err)
	}

	var payload []map[string]interface{}
	err = json.Unmarshal(jsonFile, &payload)
	if err != nil {
		log.Fatal().Msgf("Error when unmarshalling payload: ", err)
	}

	config.Elasticsearch.Hosts = []string{"http://10.1.5.14:9200"}
	// testing index
	testedIndex := "sensormap-test"

	client, err := elasticsearch.GetClient()

	for k, v := range payload {
		jsonStr, err := json.Marshal(v)
		if err != nil {
			log.Fatal().Msgf("Error when marshalling data: ", err)
		}

		req := esapi.IndexRequest{
			Index:      testedIndex,
			DocumentID: strconv.Itoa(k),
			Body:       strings.NewReader(string(jsonStr)),
			Refresh:    "true",
		}

		res, err := elasticsearch.IndexRequest(client, &req)
		if err != nil {
			log.Fatal().Msgf("Error when indexing data: ", err)
		}
		log.Info().Msgf("Success indexing data: %v", res)
	}
}
