package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

		r.GET("/person/:name", func(c *gin.Context) {
			name := c.Param("name")
			record, err := GetPersonByName(ctx, Driver, name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if record == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
				return
			}
			c.JSON(http.StatusOK, record)
		})

		r.GET("/person/region/:person_id/regions", func(c *gin.Context) {
			personID := c.Param("person_id")
			records, err := GetRegionsByPersonID(ctx, Driver, personID)
			fmt.Print(records)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if len(records) == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "Regions not found"})
				return
			}
			c.JSON(http.StatusOK, records)
		})

		r.GET("/related-persons/:personID/:degree", func(c *gin.Context) {
			personID := c.Param("personID")
			degree, err := strconv.Atoi(c.Param("degree"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid degree parameter"})
				return
			}
	
			relatedPersons, err := GetAllRelatedPersons(ctx, Driver, personID, degree)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if len(relatedPersons) == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "No related persons found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"related_persons": relatedPersons})
		})

			r.GET("/closest-relative/:personID", func(c *gin.Context) {
				personID := c.Param("personID")

				regionName := c.DefaultQuery("region", "") 
		
				if regionName == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Region name is required as a query parameter"})
					return
				}
		
				relative, err := GetClosestRelativeFromRegion(ctx, Driver, personID, regionName)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				if relative == nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "No closest relative found"})
					return
				}
				c.JSON(http.StatusOK, relative)
			})
		

	r.Run()
}
