package main

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/flavors"
	"github.com/Zillaforge/cloud-sdk/models/vps/keypairs"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
	"github.com/Zillaforge/cloud-sdk/modules/vps"
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

	// Step 1: Get default Network ID
	networksList, err := vpsClient.Networks().List(ctx, &networks.ListNetworksOptions{
		Name: "default",
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(networksList) != 1 {
		log.Fatal(errors.New("default network not found"))
	}

	defaultNetworkId := networksList[0].ID
	log.Printf("Default Network ID: %s", defaultNetworkId)

	// Step 2: Get first flavor ID from flavor list
	flavorsList, err := vpsClient.Flavors().List(ctx, &flavors.ListFlavorsOptions{})
	if err != nil {
		log.Fatal(err)
	}

	if len(flavorsList) == 0 {
		log.Fatal("No flavors found")
	}

	firstFlavorID := flavorsList[0].ID
	log.Printf("First flavor ID: %s", firstFlavorID)

	// Step 3: Create Security Group (only if it doesn't exist)
	sgName := "example-sg"

	// Check if security group already exists
	existingSGs, err := vpsClient.SecurityGroups().List(ctx, &securitygroups.ListSecurityGroupsOptions{
		Name: sgName,
	})
	if err != nil {
		log.Fatal(err)
	}

	var securityGroupID string
	if len(existingSGs) > 0 {
		// Security group already exists
		securityGroupID = existingSGs[0].ID
		log.Printf("Security group '%s' already exists with ID: %s", sgName, securityGroupID)
	} else {
		// Create new security group
		port22 := 22
		securityGroupReq := securitygroups.SecurityGroupCreateRequest{
			Name:        sgName,
			Description: "Example security group with SSH and ping access",
			Rules: []securitygroups.SecurityGroupRuleCreateRequest{
				{
					Direction:  securitygroups.DirectionIngress,
					Protocol:   securitygroups.ProtocolTCP,
					PortMin:    &port22,
					PortMax:    &port22,
					RemoteCIDR: "0.0.0.0/0",
				},
				{
					Direction:  securitygroups.DirectionIngress,
					Protocol:   securitygroups.ProtocolICMP,
					RemoteCIDR: "0.0.0.0/0",
				},
			},
		}

		createdSG, err := vpsClient.SecurityGroups().Create(ctx, securityGroupReq)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created security group: %s (ID: %s)", createdSG.Name, createdSG.ID)
	}

	// Step 4: hardcode image ID
	hardcodeImageId := os.Getenv("IMAGE_ID")
	log.Printf("Hard Code Image ID: %s", hardcodeImageId)

	// Encode password to base64
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(os.Getenv("VM_PASSWORD")))
	log.Printf("Encoded Password: %s", encodedPassword)

	key, err := vpsClient.Keypairs().Create(ctx, &keypairs.KeypairCreateRequest{
		Name: "test-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Step 5: Create Server using information from Steps 1-4
	req := &servers.ServerCreateRequest{
		Name:      "test-server",
		FlavorID:  firstFlavorID,
		ImageID:   hardcodeImageId,
		Password:  encodedPassword,
		KeypairID: key.ID,
		NICs: []servers.ServerNICCreateRequest{
			{
				NetworkID: defaultNetworkId,
				SGIDs:     []string{securityGroupID},
			},
		},
	}

	createdServer, err := vpsClient.Servers().Create(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created server: %s (ID: %s)", createdServer.Name, createdServer.ID)

	// Step 6: Wait for server to become active
	log.Printf("Waiting for server to become active...")
	err = vps.WaitForServerActive(ctx, vpsClient.Servers(), createdServer.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server is now active")

	// Step 7: Associate floating IP to server NIC
	// First, get or create a floating IP

	nicList, err := createdServer.NICs().List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var defaulNICId string

	for _, nic := range nicList {
		if nic.Network.Name == "default" {
			defaulNICId = nic.ID
		}
	}

	fipInfo, err := createdServer.NICs().AssociateFloatingIP(ctx, defaulNICId, &servers.ServerNICAssociateFloatingIPRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Associate Floating IP: %s", fipInfo.Address)

}
