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
			record, err := GetStopById(ctx, Driver, stopID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if record == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Bus stop not found"})
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
			record, err := GetStopsByRoute(ctx, Driver)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if record == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Bus stop not found"})
				return
			}
			c.JSON(http.StatusOK, record)
		})

	r.Run()
}
