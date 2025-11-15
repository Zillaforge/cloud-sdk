package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/flavors"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
	"github.com/Zillaforge/cloud-sdk/models/vps/keypairs"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
	"github.com/Zillaforge/cloud-sdk/models/vps/volumes"
	"github.com/Zillaforge/cloud-sdk/models/vrm/repositories"
	vps "github.com/Zillaforge/cloud-sdk/modules/vps/core"
	serversResource "github.com/Zillaforge/cloud-sdk/modules/vps/servers"
	vrm "github.com/Zillaforge/cloud-sdk/modules/vrm/core"
)

type App struct {
	ctx       context.Context
	vpsClient *vps.Client
	vrmClient *vrm.Client

	// Resource IDs
	defaultNetworkID  string
	firstFlavorID     string
	securityGroupID   string
	tagID             string
	keypair           *keypairs.Keypair
	server            *serversResource.ServerResource
	floatingIP        *floatingips.FloatingIP
	volume            *volumes.Volume
	repositoryID      string
	imageRepositoryID string
}

func main() {

	vpsClient, vrmClient, err := initClient()
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		ctx:       context.Background(),
		vpsClient: vpsClient,
		vrmClient: vrmClient,
	}

	// Define constants
	networkName := "default"
	securityGroupName := "example-sg"
	keypairName := "default"
	serverName := "test-server"
	volumeName := "test-vol"
	passwordEnvVar := os.Getenv("VM_PASSWORD")

	if err := app.uploadImageToRepository(); err != nil {
		log.Fatal(err)
	}

	// Setup resources
	if err := app.setupResources(networkName, securityGroupName, keypairName); err != nil {
		log.Fatal(err)
	}

	// Create server
	if err := app.createServer(serverName, passwordEnvVar); err != nil {
		log.Fatal(err)
	}

	// Create volume and attach to server
	if err := app.createVolumeAndAttach(volumeName); err != nil {
		log.Fatal(err)
	}

	// Create snapshot
	if err := app.createSnapshot(); err != nil {
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

func initClient() (*vps.Client, *vrm.Client, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Read environment variables
	baseURL := os.Getenv("API_HOST")
	protocol := os.Getenv("API_PROTOCOL")
	token := os.Getenv("API_TOKEN")
	projectCode := os.Getenv("PROJECT_SYS_CODE")

	if protocol == "" || baseURL == "" || token == "" || projectCode == "" {
		return nil, nil, errors.New("missing required environment variables: API_PROTOCOL, API_HOST, API_TOKEN, PROJECT_SYS_CODE")
	}

	ctx := context.Background()

	client, err := cloudsdk.New(protocol+"://"+baseURL, token)
	if err != nil {
		return nil, nil, err
	}

	projectClient, err := client.Project(ctx, projectCode)
	if err != nil {
		log.Fatal(err)
	}
	vpsClient := projectClient.VPS()
	vrmClient := projectClient.VRM()

	return vpsClient, vrmClient, nil
}

func (a *App) uploadImageToRepository() error {
	imageURL := "dss-public://" + os.Getenv("IMAGE_SOURCE")

	// Upload to new repository "cirros" with tag "v1"
	req := &repositories.UploadToNewRepositoryRequest{
		Name:            "cirros",
		Version:         "v1",
		Type:            "common",
		DiskFormat:      "qcow2",
		ContainerFormat: "bare",
		OperatingSystem: "linux",
		Description:     "Cirros test image",
		Filepath:        imageURL,
	}

	response, err := a.vrmClient.Repositories().Upload(a.ctx, req)
	if err != nil {
		return fmt.Errorf("failed to upload image to new repository: %w", err)
	}

	log.Printf("Successfully uploaded image to new repository: %s:%s",
		response.Repository.Name, response.Repository.Tags[0].Name)

	a.tagID = response.Repository.Tags[0].ID
	a.imageRepositoryID = response.Repository.ID

	return nil
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
	log.Printf("Retrieved default network: %s", a.defaultNetworkID)

	// Step 2: Get first flavor ID from flavor list
	flavorsList, err := a.vpsClient.Flavors().List(a.ctx, &flavors.ListFlavorsOptions{})
	if err != nil {
		return err
	}

	if len(flavorsList) == 0 {
		return errors.New("no flavors found")
	}

	a.firstFlavorID = flavorsList[0].ID
	log.Printf("Retrieved first flavor: %s", a.firstFlavorID)

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
		log.Printf("Security group '%s' already exists: %s", securityGroupName, a.securityGroupID)
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
		log.Printf("Created security group: %s (%s)", createdSG.Name, createdSG.ID)
		a.securityGroupID = createdSG.ID
	}

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
		log.Printf("Keypair '%s' already exists: %s", keypairName, a.keypair.ID)
	} else {
		// Create new keypair
		a.keypair, err = a.vpsClient.Keypairs().Create(a.ctx, &keypairs.KeypairCreateRequest{
			Name: keypairName,
		})
		if err != nil {
			return err
		}
		log.Printf("Created keypair: %s (%s)", a.keypair.Name, a.keypair.ID)
	}

	return nil
}

func (a *App) createServer(serverName, passwordEnvVar string) error {
	// Encode password to base64
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(passwordEnvVar))
	log.Printf("Encoded password: %s", encodedPassword)

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

	log.Printf("Created server: %s (%s)", a.server.Name, a.server.ID)

	// Step 7: Wait for server to become active
	log.Printf("Waiting for server to become active...")
	err = vps.WaitForServerActive(a.ctx, a.vpsClient.Servers(), a.server.ID)
	if err != nil {
		return err
	}
	log.Printf("Server is now active")

	return nil
}

func (a *App) createVolumeAndAttach(volumeName string) error {
	// Get first volume type
	volumeTypes, err := a.vpsClient.VolumeTypes().List(a.ctx)
	if err != nil {
		return err
	}

	if len(volumeTypes) == 0 {
		return errors.New("no volume types available")
	}

	firstVolumeType := volumeTypes[0]
	log.Printf("Retrieved first volume type: %s", firstVolumeType)

	// Create volume
	req := &volumes.CreateVolumeRequest{
		Name: volumeName,
		Type: firstVolumeType,
		Size: 1, // 10 GB
	}

	a.volume, err = a.vpsClient.Volumes().Create(a.ctx, req)
	if err != nil {
		return err
	}

	log.Printf("Created volume: %s (%s)", a.volume.Name, a.volume.ID)

	// Wait for volume to become available
	log.Printf("Waiting for volume to become available...")
	err = vps.WaitForVolumeAvailable(a.ctx, a.vpsClient.Volumes(), a.volume.ID)
	if err != nil {
		return err
	}
	log.Printf("Volume is now available")

	// Attach volume to server
	actionReq := &volumes.VolumeActionRequest{
		Action:   volumes.VolumeActionAttach,
		ServerID: a.server.ID,
	}

	err = a.vpsClient.Volumes().Action(a.ctx, a.volume.ID, actionReq)
	if err != nil {
		return err
	}

	log.Printf("Attached volume %s to server %s", a.volume.Name, a.server.Name)

	// Wait for volume to become in-use
	log.Printf("Waiting for volume to become in-use...")
	err = vps.WaitForVolumeInUse(a.ctx, a.vpsClient.Volumes(), a.volume.ID)
	if err != nil {
		return err
	}
	log.Printf("Volume is now in-use")

	// List all volumes attached to the server
	volumesList, err := a.server.Volumes().List(a.ctx)
	if err != nil {
		return err
	}

	for _, vol := range volumesList {
		log.Printf("Attached volume: (%s), System: %t", vol.VolumeID, vol.System)
	}

	// Detach the non-system volume
	for _, vol := range volumesList {
		if !vol.System {
			err = a.server.Volumes().Detach(a.ctx, vol.VolumeID)
			if err != nil {
				return err
			}
			log.Printf("Detached non system volume: %s", vol.Volume.Name)

			// Wait for volume to become available after detach
			log.Printf("Waiting for non system volume to become available after detach...")
			err = vps.WaitForVolumeAvailable(a.ctx, a.vpsClient.Volumes(), vol.VolumeID)
			if err != nil {
				return err
			}
			log.Printf("Non System volume is now available after detach")
		}
	}

	return nil
}

func (a *App) createSnapshot() error {
	req := &repositories.CreateSnapshotFromNewRepositoryRequest{
		Name:            "backup",
		OperatingSystem: "linux",
		Version:         time.Now().Format("2006-01-02T15:04:05"),
	}

	repoResource, err := a.vrmClient.Repositories().Snapshot(a.ctx, a.server.ID, req)
	if err != nil {
		return err
	}

	a.repositoryID = repoResource.ID
	log.Printf("Created snapshot repository: %s", a.repositoryID)

	return nil
}

func (a *App) handleFloatingIP(networkName string) error {
	// Step 8: Associate floating IP to server NIC
	nicList, err := a.server.NICs().List(a.ctx)
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

	fipInfo, err := a.server.NICs().AssociateFloatingIP(a.ctx, defaultNICId, &servers.ServerNICAssociateFloatingIPRequest{})
	if err != nil {
		return err
	}

	log.Printf("Waiting for floating IP %s to become active...", fipInfo.Address)

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

	log.Printf("Updated floating IP %s to reserved: %t", updatedFIP.Address, updatedFIP.Reserved)

	// Disassociate the floating IP
	err = a.vpsClient.FloatingIPs().Disassociate(a.ctx, a.floatingIP.ID)
	if err != nil {
		return err
	}

	log.Printf("Disassociated floating IP: %s", a.floatingIP.Address)

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
		log.Printf("Deleted floating IP: %s", a.floatingIP.Address)
	}

	// Delete Server
	if a.server != nil {
		err := a.vpsClient.Servers().Delete(a.ctx, a.server.ID)
		if err != nil {
			return err
		}
		log.Printf("Deleted server: %s", a.server.Name)
	}

	// Delete keypair
	if a.keypair != nil {
		err := a.vpsClient.Keypairs().Delete(a.ctx, a.keypair.ID)
		if err != nil {
			return err
		}
		log.Printf("Deleted keypair: %s", a.keypair.Name)
	}

	// Delete SG
	if a.securityGroupID != "" {
		err := a.vpsClient.SecurityGroups().Delete(a.ctx, a.securityGroupID)
		if err != nil {
			return err
		}
		log.Printf("Deleted security group: %s", "example-sg")
	}

	// Wait for volume to be available, then delete
	if a.volume != nil {
		err := a.vpsClient.Volumes().Delete(a.ctx, a.volume.ID)
		if err != nil {
			return err
		}
		log.Printf("Deleted volume: %s", a.volume.Name)
	}

	// Delete repositories
	for _, repo := range []string{a.repositoryID, a.imageRepositoryID} {
		if repo != "" {
			err := a.vrmClient.Repositories().Delete(a.ctx, repo)
			if err != nil {
				return err
			}
			log.Printf("Deleted repoistory: %s", repo)
		}
	}

	return nil
}
