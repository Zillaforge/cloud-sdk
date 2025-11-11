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
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
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
		securityGroupID = createdSG.ID
	}

	// Step 4: hardcode image ID
	hardcodeImageId := os.Getenv("IMAGE_ID")
	log.Printf("Hard Code Image ID: %s", hardcodeImageId)

	// Encode password to base64
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(os.Getenv("VM_PASSWORD")))
	log.Printf("Encoded Password: %s", encodedPassword)

	// Step 5: Create keypair if not exist
	keypairName := "default"
	existingKeypairs, err := vpsClient.Keypairs().List(ctx, &keypairs.ListKeypairsOptions{
		Name: keypairName,
	})
	if err != nil {
		log.Fatal(err)
	}

	var key *keypairs.Keypair
	if len(existingKeypairs) > 0 {
		// Keypair already exists
		key = existingKeypairs[0]
		log.Printf("Keypair '%s' already exists with ID: %s", keypairName, key.ID)
	} else {
		// Create new keypair
		key, err = vpsClient.Keypairs().Create(ctx, &keypairs.KeypairCreateRequest{
			Name: keypairName,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created keypair: %s (ID: %s)", key.Name, key.ID)
	}

	// Step 6: Create Server using information from Steps 1-4
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

	// Step 7: Wait for server to become active
	log.Printf("Waiting for server to become active...")
	err = vps.WaitForServerActive(ctx, vpsClient.Servers(), createdServer.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server is now active")

	// Step 8: Associate floating IP to server NIC

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

	log.Printf("Wait for Associate Floating IP become ACTIVE: %s", fipInfo.Address)

	err = vps.WaitForFloatingIPActive(ctx, vpsClient.FloatingIPs(), fipInfo.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Step 9: Update associated Floating IP to Reserved, then Disassociate

	associateIP, err := vpsClient.FloatingIPs().List(ctx, &floatingips.ListFloatingIPsOptions{
		Address: fipInfo.Address,
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(associateIP) == 0 {
		log.Fatalf("Floating IP %s not found", fipInfo.Address)
	}

	// Update the floating IP to Reserved
	reserved := true
	updateReq := &floatingips.FloatingIPUpdateRequest{
		Reserved: &reserved,
	}

	updatedFIP, err := vpsClient.FloatingIPs().Update(ctx, associateIP[0].ID, updateReq)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Updated Floating IP %s to Reserved: %t", updatedFIP.Address, updatedFIP.Reserved)

	// Disassociate the floating IP
	err = vpsClient.FloatingIPs().Disassociate(ctx, associateIP[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Disassociated Floating IP: %s", associateIP[0].Address)

	// Step 10: teardown

	// Delete FIP
	err = vpsClient.FloatingIPs().Delete(ctx, associateIP[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Floating IP %s Deleted", associateIP[0].Address)

	// Delete Server
	err = vpsClient.Servers().Delete(ctx, createdServer.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server %s Deleted", createdServer.Name)

	// Delete keypair
	err = vpsClient.Keypairs().Delete(ctx, key.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Keypair %s deleted", key.Name)

	// Delete SG
	err = vpsClient.SecurityGroups().Delete(ctx, securityGroupID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Security group %s deleted", "default")

}
