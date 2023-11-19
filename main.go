package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	var err error
	Driver, err = InitDriver("neo4j://localhost:7687", "neo4j", "password")
	if err != nil {
		log.Fatalf("Failed to create the driver: %v", err)
	}
	defer Driver.Close(ctx)

	r := gin.Default()

	r.GET("/stop/:id", func(c *gin.Context) {
		stopID := c.Param("id")
		record, err := getStopById(ctx, Driver, stopID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if record == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Buss stop not found"})
			return
		}
		c.JSON(http.StatusOK, record)
	})

	r.GET("/stop/:id/routes", func(c *gin.Context) {
		stopID := c.Param("id")
		records, err := getStopRoutes(ctx, Driver, stopID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(records) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Routes not found"})
			return
		}
		c.JSON(http.StatusOK, records)
	})

	r.GET("/route/:id/stops", func(c *gin.Context) {
		routeID := c.Param("id")
		records, err := getStopsByRoute(ctx, Driver, routeID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(records) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stops not found"})
			return
		}
		c.JSON(http.StatusOK, records)
	})

	r.GET("/all_routes/:start_stop_id/:end_stop_id", func(c *gin.Context) {
		startStopID := c.Param("start_stop_id")
		endStopID := c.Param("end_stop_id")
		records, err := getAllRoutesBetweenStops(ctx, Driver, startStopID, endStopID)
		fmt.Print(records)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(records) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Routes not found"})
			return
		}
		c.JSON(http.StatusOK, records)
	})

	r.GET("/optimal_route/:start_stop_id/:end_stop_id", func(c *gin.Context) {
		startStopID := c.Param("start_stop_id")
		endStopID := c.Param("end_stop_id")
		records, err := getShortestPathByBus(ctx, Driver, startStopID, endStopID)
		fmt.Print(records)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(records) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
			return
		}
		c.JSON(http.StatusOK, records)
	})

	r.GET("/stops_by_route", func(c *gin.Context) {
		record, err := getStopsCountByRoute(ctx, Driver)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if record == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stops not found"})
			return
		}
		c.JSON(http.StatusOK, record)
	})

	r.GET("/route/:id/total_distance", func(c *gin.Context) {
		routeID := c.Param("id")
		routeInfo, err := calculateTotalDistance(ctx, Driver, routeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if routeInfo == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Route stop not found"})
			return
		}
		c.JSON(http.StatusOK, routeInfo)
	})

	r.Run()
}
