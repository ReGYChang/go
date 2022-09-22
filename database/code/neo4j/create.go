package main

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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

func insertSn(session neo4j.Session, sn, kpsn, moNumber string, workTime time.Time) error {
	params := map[string]interface{}{
		"master_sn": sn,
		"slave_sn":  kpsn,
		"mo_number": moNumber,
		"work_time": workTime,
	}
	_, err := session.Run(`
		MERGE (mo: MO { moNumber: $mo_number})
		ON CREATE
			SET mo.mo = $mo_number
		MERGE (master { serialNumber: $master_sn, workTime: $work_time })
		ON CREATE
			SET master: $serial_number_type
		ON MATCH
			SET master: $serial_number_type
		MERGE (slave { serialNumber: $slave_sn, workTime: $work_time })
		ON CREATE
			SET slave: $serial_number_type
		ON MATCH
			SET slave: $serial_number_type
		MERGE (slave)-[: BELONGS]->(master)
		MERGE (slave)-[: BELONGS]->(mo)
		MERGE (master)-[: BELONGS]->(mo)
	`, params)

	if err != nil {
		return err
	}

	return nil
}
