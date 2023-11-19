package main

import (
	"context"

	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)


func GetStopById(ctx context.Context, driver neo4j.DriverWithContext, stopID string) (*neo4j.Record, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := "MATCH (s:Stop {stop_id: $stopID}) RETURN s"
	parameters := map[string]any{
		"stopID": stopID,
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


func getStopRoutes(ctx context.Context, driver neo4j.DriverWithContext, stopID string) ([]*neo4j.Record, error)  {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

    query := "MATCH(s:Stop {stop_id: $stopID})-[:SERVICED_BY]->(r:Route) RETURN r"
    parameters := map[string]any{
		"stopID": stopID,
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

func getAllRoutesBetweenStops(ctx context.Context, driver neo4j.DriverWithContext, startStopID string, endStopID string) ([]*neo4j.Record, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
    MATCH (start:Stop {stop_id: $startStopID}), (end:Stop {stop_id: $endStopID})
    MATCH path = (start)-[:SERVICED_BY|SEGMENT*]-(end)
    WITH [r in nodes(path) WHERE 'Route' in labels(r)] as routes
    UNWIND routes as route
    RETURN DISTINCT route.name as RouteName    
	`
	parameters := map[string]any{
		"startStopID": startStopID,
		"endStopID":   endStopID,
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

func getShortestPathByBus(ctx context.Context, driver neo4j.DriverWithContext, startStopID string, endStopID string) ([]*neo4j.Record, error)  {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
	MATCH (start:Stop {stop_id: $startStopID}), (end:Stop {stop_id: $endStopID})
	CALL apoc.algo.dijkstra(start, end, 'SEGMENT>', 'distance') 
	YIELD path, weight
	RETURN nodes(path) AS nodePath, weight AS totalDistance 
    `
	parameters := map[string]any{
		"startStopID": startStopID,
		"endStopID":   endStopID,
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
    fmt.Print(result)

	if err := result.Err(); err != nil {
		return nil, err
	}

	return records, nil
}


func GetStopsByRoute(ctx context.Context, driver neo4j.DriverWithContext)  ([]*neo4j.Record, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `MATCH (r:Route)-[:SERVICED_BY]-(s:Stop)
    RETURN r.name AS Route, count(DISTINCT s) AS NumberOfStops`

    result, err := session.Run(ctx, query, nil)
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