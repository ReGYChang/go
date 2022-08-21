package main

import (
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
	testSn := "2M1829A8NEFR"

	if err := findRoot(session, testSn); err != nil {
		panic(err)
	}
}

func findRoot(session neo4j.Session, sn string) error {
	params := map[string]interface{}{
		"sn": sn,
	}
	_, err := session.Run(`
		MATCH (n)
		WHERE n.sn = $sn
		WITH n, [(n) - [:BELONGS*] -> (x:MO) | x] AS mo
		RETURN n, mo
	`, params)

	if err != nil {
		return err
	}

	return nil
}
