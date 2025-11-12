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
	"github.com/Zillaforge/cloud-sdk/models/vrm/tags"
	"github.com/Zillaforge/cloud-sdk/modules/vps"
	serversResource "github.com/Zillaforge/cloud-sdk/modules/vps/servers"
	"github.com/Zillaforge/cloud-sdk/modules/vrm"
)

type App struct {
	ctx       context.Context
	vpsClient *vps.Client
	vrmClient *vrm.Client

	// Resource IDs
	defaultNetworkID string
	firstFlavorID    string
	securityGroupID  string
	tagID            string
	keypair          *keypairs.Keypair
	server           interface{} // *servers.ServerResource
	floatingIP       *floatingips.FloatingIP
}

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
	vrmClient := client.Project(projectID).VRM()

	app := &App{
		ctx:       ctx,
		vpsClient: vpsClient,
		vrmClient: vrmClient,
	}

	// Define constants
	networkName := "default"
	securityGroupName := "example-sg"
	keypairName := "default"
	serverName := "test-server"
	passwordEnvVar := os.Getenv("VM_PASSWORD")

	// Setup resources
	if err := app.setupResources(networkName, securityGroupName, keypairName); err != nil {
		log.Fatal(err)
	}

	// Create server
	if err := app.createServer(serverName, passwordEnvVar); err != nil {
		log.Fatal(err)
	}

	// Handle floating IP
	if err := app.handleFloatingIP(networkName); err != nil {
		log.Fatal(err)
	}

	// Teardown resources
	if err := app.teardown(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) setupResources(networkName, securityGroupName, keypairName string) error {
	// Step 1: Get default Network ID
	networksList, err := a.vpsClient.Networks().List(a.ctx, &networks.ListNetworksOptions{
		Name: networkName,
	})
	if err != nil {
		return err
	}

	if len(networksList) != 1 {
		return errors.New("default network not found")
	}

	a.defaultNetworkID = networksList[0].ID
	log.Printf("Default Network ID: %s", a.defaultNetworkID)

	// Step 2: Get first flavor ID from flavor list
	flavorsList, err := a.vpsClient.Flavors().List(a.ctx, &flavors.ListFlavorsOptions{})
	if err != nil {
		return err
	}

	if len(flavorsList) == 0 {
		return errors.New("no flavors found")
	}

	a.firstFlavorID = flavorsList[0].ID
	log.Printf("First flavor ID: %s", a.firstFlavorID)

	// Step 3: Create Security Group (only if it doesn't exist)
	// Check if security group already exists
	existingSGs, err := a.vpsClient.SecurityGroups().List(a.ctx, &securitygroups.ListSecurityGroupsOptions{
		Name: securityGroupName,
	})
	if err != nil {
		return err
	}

	if len(existingSGs) > 0 {
		// Security group already exists
		a.securityGroupID = existingSGs[0].ID
		log.Printf("Security group '%s' already exists with ID: %s", securityGroupName, a.securityGroupID)
	} else {
		// Create new security group
		port22 := 22
		securityGroupReq := securitygroups.SecurityGroupCreateRequest{
			Name:        securityGroupName,
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

		createdSG, err := a.vpsClient.SecurityGroups().Create(a.ctx, securityGroupReq)
		if err != nil {
			return err
		}
		log.Printf("Created security group: %s (ID: %s)", createdSG.Name, createdSG.ID)
		a.securityGroupID = createdSG.ID
	}

	// Step 4: Get first tag ID
	tagsList, err := a.vrmClient.Tags().List(a.ctx, &tags.ListTagsOptions{})
	if err != nil {
		return err
	}

	if len(tagsList) == 0 {
		return errors.New("no tag found")
	}
	a.tagID = tagsList[0].ID
	log.Printf("First Repository/Tag is %s:%s", tagsList[0].Repository.Name, tagsList[0].Name)

	// Step 5: Create keypair if not exist
	existingKeypairs, err := a.vpsClient.Keypairs().List(a.ctx, &keypairs.ListKeypairsOptions{
		Name: keypairName,
	})
	if err != nil {
		return err
	}

	if len(existingKeypairs) > 0 {
		// Keypair already exists
		a.keypair = existingKeypairs[0]
		log.Printf("Keypair '%s' already exists with ID: %s", keypairName, a.keypair.ID)
	} else {
		// Create new keypair
		a.keypair, err = a.vpsClient.Keypairs().Create(a.ctx, &keypairs.KeypairCreateRequest{
			Name: keypairName,
		})
		if err != nil {
			return err
		}
		log.Printf("Created keypair: %s (ID: %s)", a.keypair.Name, a.keypair.ID)
	}

	return nil
}

func (a *App) createServer(serverName, passwordEnvVar string) error {
	// Encode password to base64
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(os.Getenv(passwordEnvVar)))
	log.Printf("Encoded Password: %s", encodedPassword)

	// Step 6: Create Server using information from setup
	req := &servers.ServerCreateRequest{
		Name:      serverName,
		FlavorID:  a.firstFlavorID,
		ImageID:   a.tagID,
		Password:  encodedPassword,
		KeypairID: a.keypair.ID,
		NICs: []servers.ServerNICCreateRequest{
			{
				NetworkID: a.defaultNetworkID,
				SGIDs:     []string{a.securityGroupID},
			},
		},
	}

	var err error
	a.server, err = a.vpsClient.Servers().Create(a.ctx, req)
	if err != nil {
		return err
	}

	server := a.server.(*serversResource.ServerResource)
	log.Printf("Created server: %s (ID: %s)", server.Name, server.ID)

	// Step 7: Wait for server to become active
	log.Printf("Waiting for server to become active...")
	err = vps.WaitForServerActive(a.ctx, a.vpsClient.Servers(), server.ID)
	if err != nil {
		return err
	}
	log.Printf("Server is now active")

	return nil
}

func (a *App) handleFloatingIP(networkName string) error {
	// Step 8: Associate floating IP to server NIC
	server := a.server.(*serversResource.ServerResource)
	nicList, err := server.NICs().List(a.ctx)
	if err != nil {
		return err
	}

	var defaultNICId string
	for _, nic := range nicList {
		if nic.Network.Name == networkName {
			defaultNICId = nic.ID
			break
		}
	}

	fipInfo, err := server.NICs().AssociateFloatingIP(a.ctx, defaultNICId, &servers.ServerNICAssociateFloatingIPRequest{})
	if err != nil {
		return err
	}

	log.Printf("Wait for Associate Floating IP become ACTIVE: %s", fipInfo.Address)

	err = vps.WaitForFloatingIPActive(a.ctx, a.vpsClient.FloatingIPs(), fipInfo.ID)
	if err != nil {
		return err
	}

	// Step 9: Update associated Floating IP to Reserved, then Disassociate
	associateIP, err := a.vpsClient.FloatingIPs().List(a.ctx, &floatingips.ListFloatingIPsOptions{
		Address: fipInfo.Address,
	})
	if err != nil {
		return err
	}

	if len(associateIP) == 0 {
		return errors.New("floating IP not found")
	}

	a.floatingIP = associateIP[0]

	// Update the floating IP to Reserved
	reserved := true
	updateReq := &floatingips.FloatingIPUpdateRequest{
		Reserved: &reserved,
	}

	updatedFIP, err := a.vpsClient.FloatingIPs().Update(a.ctx, a.floatingIP.ID, updateReq)
	if err != nil {
		return err
	}

	log.Printf("Updated Floating IP %s to Reserved: %t", updatedFIP.Address, updatedFIP.Reserved)

	// Disassociate the floating IP
	err = a.vpsClient.FloatingIPs().Disassociate(a.ctx, a.floatingIP.ID)
	if err != nil {
		return err
	}

	log.Printf("Disassociated Floating IP: %s", a.floatingIP.Address)

	return nil
}

func (a *App) teardown() error {
	// Step 10: Teardown resources

	// Delete FIP
	if a.floatingIP != nil {
		err := a.vpsClient.FloatingIPs().Delete(a.ctx, a.floatingIP.ID)
		if err != nil {
			return err
		}
		log.Printf("Floating IP %s Deleted", a.floatingIP.Address)
	}

	// Delete Server
	if a.server != nil {
		server := a.server.(*serversResource.ServerResource)
		err := a.vpsClient.Servers().Delete(a.ctx, server.ID)
		if err != nil {
			return err
		}
		log.Printf("Server %s Deleted", server.Name)
	}

	// Delete keypair
	if a.keypair != nil {
		err := a.vpsClient.Keypairs().Delete(a.ctx, a.keypair.ID)
		if err != nil {
			return err
		}
		log.Printf("Keypair %s deleted", a.keypair.Name)
	}

	// Delete SG
	if a.securityGroupID != "" {
		err := a.vpsClient.SecurityGroups().Delete(a.ctx, a.securityGroupID)
		if err != nil {
			return err
		}
		log.Printf("Security group %s deleted", "example-sg")
	}

	return nil
}
