package main

import (
	"context"

	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func getStopById(ctx context.Context, driver neo4j.DriverWithContext, stopID string) (map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := "MATCH (s:Stop {stop_id: $stopID}) RETURN s.stop_id AS stop_id, s.name AS name"
	parameters := map[string]interface{}{
		"stopID": stopID,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	var stop map[string]interface{}

	if result.Next(ctx) {
		record := result.Record()
		stopID, ok := record.Get("stop_id")
		if !ok {
			return nil, fmt.Errorf("stop_id not found in the record")
		}
		name, ok := record.Get("name")
		if !ok {
			return nil, fmt.Errorf("name not found in the record")
		}
		stop = map[string]interface{}{
			"stop_id": stopID,
			"name":    name,
		}
	} else {
		if err := result.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("no stop found with stop_id: %s", stopID)
	}

	return stop, nil
}

func getStopRoutes(ctx context.Context, driver neo4j.DriverWithContext, stopID string) ([]map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := "MATCH (s:Stop {stop_id: $stopID})-[:SERVICED_BY]->(r:Route) RETURN r.name AS name, r.route_id AS route_id"
	parameters := map[string]any{
		"stopID": stopID,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	var routes []map[string]interface{}
	for result.Next(ctx) {
		record := result.Record()
		name, _ := record.Get("name")
		routeID, _ := record.Get("route_id")
		route := map[string]interface{}{
			"name":     name,
			"route_id": routeID,
		}
		routes = append(routes, route)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func getStopsByRoute(ctx context.Context, driver neo4j.DriverWithContext, routeID string) ([]map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
	MATCH (route:Route {route_id: $routeID})<-[:SERVICED_BY]-(stop:Stop)
	RETURN stop.name AS name, stop.stop_id AS stop_id
	`
	parameters := map[string]any{
		"routeID": routeID,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	var stops []map[string]interface{}
	for result.Next(ctx) {
		record := result.Record()
		stopName, _ := record.Get("name")
		stopID, _ := record.Get("stop_id")
		stops = append(stops, map[string]interface{}{
			"name":    stopName,
			"stop_id": stopID,
		})
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return stops, nil
}

func getAllRoutesBetweenStops(ctx context.Context, driver neo4j.DriverWithContext, startStopID string, endStopID string) ([]map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
	MATCH (start:Stop {stop_id: $startStopID})-[:SERVICED_BY]->(route:Route)<-[:SERVICED_BY]-(end:Stop {stop_id: $endStopID})
	RETURN DISTINCT route.name AS name
	`

	parameters := map[string]any{
		"startStopID": startStopID,
		"endStopID":   endStopID,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	var routes []map[string]interface{}
	for result.Next(ctx) {
		record := result.Record()
		name, _ := record.Get("name")
		route := map[string]interface{}{
			"name": name,
		}
		routes = append(routes, route)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func getShortestPathByBus(ctx context.Context, driver neo4j.DriverWithContext, startStopID string, endStopID string) (map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
    MATCH (start:Stop {stop_id: $startStopID}), (end:Stop {stop_id: $endStopID})
	WHERE start <> end
    CALL apoc.algo.dijkstra(start, end, 'SEGMENT', 'distance') 
    YIELD path, weight
    UNWIND nodes(path) AS n
    RETURN weight AS totalDistance, COLLECT({name: n.name, stop_id: n.stop_id}) AS Stops
    `
	parameters := map[string]any{
		"startStopID": startStopID,
		"endStopID":   endStopID,
	}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		totalDistance, _ := record.Get("totalDistance")
		stops, _ := record.Get("Stops")

		pathInfo := map[string]interface{}{
			"Route": map[string]interface{}{
				"distance": totalDistance,
				"Stops":    stops,
			},
		}
		return pathInfo, nil
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("no path found between stop IDs %s and %s", startStopID, endStopID)
}

func getStopsCountByRoute(ctx context.Context, driver neo4j.DriverWithContext) ([]map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
	MATCH (r:Route)-[:SERVICED_BY]-(s:Stop)
	RETURN r.name AS name, COUNT(DISTINCT s) AS stops
	`

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var routes []map[string]interface{}
	for result.Next(ctx) {
		record := result.Record()
		routeName, _ := record.Get("name")
		stopsCount, _ := record.Get("stops")
		route := map[string]interface{}{
			"name":  routeName,
			"stops": stopsCount,
		}
		routes = append(routes, route)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func calculateTotalDistance(ctx context.Context, driver neo4j.DriverWithContext, routeID string) (map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (route:Route {route_id: $routeID})
		MATCH (stop1:Stop)-[seg:SEGMENT]->(stop2:Stop)
		WHERE (stop1)-[:SERVICED_BY]->(route) AND (stop2)-[:SERVICED_BY]->(route)
		RETURN route.name AS routeName, SUM(seg.distance) AS totalDistance
    `
	parameters := map[string]interface{}{"routeID": routeID}

	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		totalDistance, ok := record.Get("totalDistance")
		if !ok {
			return nil, fmt.Errorf("totalDistance not found in the record")
		}
		routeName, ok := record.Get("routeName")
		if !ok {
			return nil, fmt.Errorf("routeName not found in the record")
		}
		return map[string]interface{}{
			"route_name":     routeName,
			"total_distance": totalDistance,
		}, nil
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("no route found with route_id: %s", routeID)
}
