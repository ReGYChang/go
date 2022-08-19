package main

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

func main() {
	dbUri := "neo4j://localhost:7687"
	driver, err := neo4j.NewDriver(dbUri, neo4j.BasicAuth("neo4j", "regy", ""))
	if err != nil {
		panic(err)
	}

	// Handle driver lifetime based on your application lifetime requirements  driver's lifetime is usually
	// bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	testSet := [][]string{[]string{"1", "2"}, []string{"1", "3"}, []string{"2", "4"}}

	for _, v := range testSet {
		if err := insertSn(session, v[0], v[1]); err != nil {
			panic(err)
		}
	}
}

func insertSn(session neo4j.Session, sn string, kpsn string) error {
	params := map[string]interface{}{
		"sn":   sn,
		"kpsn": kpsn,
	}
	_, err := session.Run(`
		MERGE (sn { sn: $sn })
		ON CREATE
			SET sn: SN
		ON MATCH
			SET sn: SN
		MERGE (kpsn { sn: $kpsn })
		ON CREATE
			SET kpsn: KPSN
		ON MATCH
			SET kpsn: KPSN
		MERGE (kpsn)-[r: BELONGS]->(sn)
	`, params)

	if err != nil {
		return err
	}

	return nil
}