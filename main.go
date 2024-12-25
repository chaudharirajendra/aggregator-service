package main

import (
	clients "feed-service/client"
	"feed-service/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Fetch service URLs from environment variables
	repairServiceURL := os.Getenv("REPAIR_SERVICE_URL")
	sessionServiceURL := os.Getenv("SESSION_SERVICE_URL")
	machineNameServiceURL := os.Getenv("MACHINE_SERVICE_URL")

	// Validate the environment variables
	if repairServiceURL == "" || sessionServiceURL == "" || machineNameServiceURL == "" {
		log.Fatal("Service URLs are not set properly. Please set REPAIR_SERVICE_URL, SESSION_SERVICE_URL, and MACHINE_NAME_SERVICE_URL environment variables.")
	}

	client := clients.NewHTTPMachineServiceClient(repairServiceURL, sessionServiceURL, machineNameServiceURL)
	handler := handlers.NewMachineFeedHandler(client)

	router := gin.Default()
	router.GET("/machine-feeds/:machineId", handler.GetFeeds)

	port := 8080
	println("Feed service running at http://localhost:", port)
	router.Run(":8080")
}
