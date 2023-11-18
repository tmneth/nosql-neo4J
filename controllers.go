package main

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)


func GetPersonByName(ctx context.Context, driver neo4j.DriverWithContext, name string) (*neo4j.Record, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := "MATCH (p:Person {name: $name}) RETURN p"
	parameters := map[string]any{
		"name": name,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		return record, nil
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func GetRegionsByPersonID(ctx context.Context, driver neo4j.DriverWithContext, personID string) ([]*neo4j.Record, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := "MATCH (p:Person {person_id: $personID})-[:FROM_REGION]->(r:Region) RETURN r"
	parameters := map[string]any{
		"personID": personID,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	var records []*neo4j.Record
	for result.Next(ctx) {
		record := result.Record()
		records = append(records, record)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func GetAllRelatedPersons(ctx context.Context, driver neo4j.DriverWithContext, personID string, degree int) ([]map[string]interface{}, error) {
    session := driver.NewSession(ctx, neo4j.SessionConfig{})
    defer session.Close(ctx)

    query := `
    MATCH path = (startPerson:Person {person_id: $personID})-[:RELATED_TO*]-(relatedPerson:Person)
    WHERE startPerson <> relatedPerson AND length(path) <= $degree
    RETURN relatedPerson AS person, length(path) AS degree
    `
    parameters := map[string]interface{}{
        "personID": personID,
        "degree":   degree,
    }

    result, err := session.Run(ctx, query, parameters)
    if err != nil {
        return nil, err
    }

    var relatedPersons []map[string]interface{}
    for result.Next(ctx) {
        record := result.Record()
        personNode, _ := record.Get("person")
        person, _ := personNode.(neo4j.Node)
        relDegree, _ := record.Get("degree")

        relatedPersons = append(relatedPersons, map[string]interface{}{
            "related_person": map[string]interface{}{
                "person_id": person.Props["person_id"],
                "name":      person.Props["name"],
            },
            "degree": relDegree,
        })
    }

    if err := result.Err(); err != nil {
        return nil, err
    }

    return relatedPersons, nil
}


func GetClosestRelativeFromRegion(ctx context.Context, driver neo4j.DriverWithContext, personID, regionName string) (map[string]interface{}, error) {
    session := driver.NewSession(ctx, neo4j.SessionConfig{})
    defer session.Close(ctx)

    query := `
	MATCH (person:Person {person_id: $personID}),
	(region:Region {name: $regionName}),
	path = (person)-[:RELATED_TO*]-(relative:Person)
	WHERE (relative)-[:FROM_REGION]->(region)
	RETURN relative, length(path) as degreesOfSeparation
	ORDER BY degreesOfSeparation ASC
	LIMIT 1

    `
    parameters := map[string]interface{}{
        "personID":   personID,
        "regionName": regionName,
    }

    result, err := session.Run(ctx, query, parameters)
    if err != nil {
        return nil, err
    }

    if result.Next(ctx) {
        record := result.Record()
        relativeNode, _ := record.Get("relative")
        relative, _ := relativeNode.(neo4j.Node)
        degreesOfSeparation, _ := record.Get("degreesOfSeparation")

        return map[string]interface{}{
            "relative": map[string]interface{}{
                "person_id": relative.Props["person_id"],
                "name":      relative.Props["name"],
            },
            "degreesOfSeparation": degreesOfSeparation,
        }, nil
    }

    if err := result.Err(); err != nil {
        return nil, err
    }

    return nil, nil
}