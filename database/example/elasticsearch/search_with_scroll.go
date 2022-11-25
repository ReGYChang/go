package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/buger/jsonparser"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

var countSuccessful uint64

func main() {
	config.Elasticsearch.Hosts = []string{"http://10.13.1.27:9200"}
	client, err := GetClient()
	if err != nil {
		panic(err)
	}

	var batchNum int

	reqBody := fmt.Sprintf(`
		{
		  "query": {
			"match_all": {}
		  }
		}
	`)

	//scrolling
	res, _ := client.Search(
		client.Search.WithIndex(fmt.Sprintf("index")),
		client.Search.WithBody(strings.NewReader(reqBody)),
		client.Search.WithSize(10000),
		client.Search.WithSort("_doc"),
		client.Search.WithScroll(5*time.Second),
	)

	// Handle the first batch of data and extract the scrollID
	//
	data := read(res.Body)
	res.Body.Close()

	scrollID := gjson.Get(data, "_scroll_id").String()

	fmt.Println("Batch   ", batchNum)
	fmt.Println("ScrollID", scrollID)
	fmt.Println("IDs     ", gjson.Get(data, "hits.hits.#._id"))
	fmt.Println(strings.Repeat("-", 80))

	for {
		batchNum++

		// Perform the scroll request and pass the scrollID and scroll duration
		//
		res, err := client.Scroll(client.Scroll.WithScrollID(scrollID), client.Scroll.WithScroll(5*time.Second))
		if err != nil {
			panic(err)
		}
		if res.IsError() {
			panic(res)
		}

		data := read(res.Body)
		res.Body.Close()

		// Extract the scrollID from response
		//
		scrollID = gjson.Get(data, "_scroll_id").String()

		// Extract the search results
		//
		hits := gjson.Get(data, "hits.hits")

		for _, d := range hits.Array() {
			b, _, _, _ := jsonparser.Get([]byte(d.String()), "_source")
			atomic.AddUint64(&countSuccessful, 1)
			log.Info().Str("Scroll", strconv.FormatUint(countSuccessful, 10)).Msg(string(b))
		}

		// Break out of the loop when there are no results
		//
		if len(hits.Array()) < 1 {
			fmt.Println("Finished scrolling")
			break
		} else {
			fmt.Println("Batch   ", batchNum)
			fmt.Println("ScrollID", scrollID)
			fmt.Println("IDs     ", gjson.Get(hits.Raw, "#._id"))
			fmt.Println(strings.Repeat("-", 80))
		}
	}
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}
