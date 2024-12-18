package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mrbryside/config"
	"github.com/mrbryside/github"
	"log"
	"net/http"
)

func main() {
	// init config
	config.LoadConfig()

	// start echo
	e := echo.New()
	e.POST("/pr/review", prHandler)

	log.Fatal(e.Start(":8080"))
}

type requestBody struct {
	Owner          string `json:"owner"`
	RepositoryName string `json:"repository_name"`
	PrNumber       int    `json:"pr_number"`
}

func prHandler(c echo.Context) error {
	req := new(requestBody)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request",
		})
	}

	changes, err := github.FetchPRChanges(req.Owner, req.RepositoryName, req.PrNumber)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to fetch PR changes",
		})
	}

	go func() {
		// trigger ai reviewer
		fmt.Println(changes)
	}()

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully triggered ai reviewer!",
	})
}
