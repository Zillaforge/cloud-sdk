package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read environment variables
	baseURL := os.Getenv("API_HOST")
	protocol := os.Getenv("API_PROTOCOL")
	token := os.Getenv("API_TOKEN")
	projectID := os.Getenv("PROJECT_ID")

	if protocol == "" || baseURL == "" || token == "" || projectID == "" {
		log.Fatal("Missing required environment variables: BASE_URL, TOKEN, PROJECT_ID")
	}

	ctx := context.Background()

	client, err := cloudsdk.New(protocol+"://"+baseURL, token)
	if err != nil {
		log.Fatal(err)
	}

	vpsClient := client.Project(projectID).VPS()

	networksList, err := vpsClient.Networks().List(ctx, &networks.ListNetworksOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range networksList.Networks {
		log.Println(n.Name)
	}

	if len(networksList.Networks) != 1 {
		log.Fatalf("expected 1 network, got %d", len(networksList.Networks))
	}
}
