package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"nexdata/pkg/config"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func TestGetClient(t *testing.T) {
	// 1. Set up testing configuration
	//
	config.Elasticsearch.Hosts = []string{"http://10.1.5.14:9200"}
	// testing index
	testedIndex := "test-2022.07.14"

	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)

	// 2. Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	es, err := GetClient()
	if err != nil {
		t.Errorf("Test get elasticsearch client failed: %v", err)
	}

	// 3. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		t.Errorf("Test get elasticsearch cluster info failed: %s", err)
	}
	// Check response status
	if res.IsError() {
		t.Errorf("Error getting cluster info response: %s", res.Status())
	}
	// Deserialize the response into a map
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Errorf("Error parsing the reponse body: %s", err)
	}
	t.Logf("Client: %s", elasticsearch.Version)
	t.Logf("Server: %s", r["version"].(map[string]interface{})["number"])
	t.Log(strings.Repeat("~", 37))

	// 4. Index documents concurrently
	//
	for i, title := range []string{"Test One", "Test Two"} {
		wg.Add(1)

		go func(i int, title string) {
			defer wg.Done()

			// Build the request body.
			var b strings.Builder
			b.WriteString(`{"title" : "`)
			b.WriteString(title)
			b.WriteString(`"}`)

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      testedIndex,
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(b.String()),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				t.Errorf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				t.Errorf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					t.Logf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					t.Logf("response status: [%s] %s; indexed document version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, title)
	}
	wg.Wait()

	t.Log(strings.Repeat("-", 37))

	// 5. Search for the indexed documents
	//
	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		t.Errorf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(testedIndex),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		t.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			t.Errorf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			t.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Errorf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	t.Logf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		t.Logf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	t.Log(strings.Repeat("=", 37))
}
