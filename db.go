package main

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Driver neo4j.DriverWithContext

func InitDriver(uri, username, password string) (neo4j.DriverWithContext, error) {
	var err error
	Driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	return Driver, nil
}
