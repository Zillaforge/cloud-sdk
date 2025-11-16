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
	"github.com/Zillaforge/cloud-sdk/models/vps/snapshots"
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

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Define constants
	networkName := "default"
	securityGroupName := "default-sg"
	keypairName := "default"
	serverName := "default"
	volumeName := "default"

	baseURL := os.Getenv("API_HOST")
	protocol := os.Getenv("API_PROTOCOL")
	token := os.Getenv("API_TOKEN")
	projectCode := os.Getenv("PROJECT_SYS_CODE")
	passwordEnvVar := os.Getenv("VM_PASSWORD")
	imageURL := "dss-public://" + os.Getenv("IMAGE_SOURCE")

	// 1. 初始化客戶端
	vpsClient, vrmClient, err := initClient(protocol, baseURL, token, projectCode)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		ctx:       context.Background(),
		vpsClient: vpsClient,
		vrmClient: vrmClient,
	}

	// 2. 取得Network與Flavor
	if err := app.getNetworkAndFlavor(); err != nil {
		log.Fatal(err)
	}

	// 3. 上傳Image到倉庫
	if err := app.uploadImageToRepository(imageURL); err != nil {
		log.Fatal(err)
	}

	// 4. 建立或檢查安全群組
	if err := app.createOrCheckSecurityGroup(securityGroupName); err != nil {
		log.Fatal(err)
	}

	// 5. 建立或檢查 Keypair
	if err := app.createOrCheckKeypair(keypairName); err != nil {
		log.Fatal(err)
	}

	// 6. 建立伺服器
	if err := app.createServer(serverName, passwordEnvVar); err != nil {
		log.Fatal(err)
	}

	// 7. 建立 Volume
	if err := app.createVolume(volumeName); err != nil {
		log.Fatal(err)
	}

	// 8. 建立 Volume Snapshot
	if err := app.createVolumeSnapshot(); err != nil {
		log.Fatal(err)
	}

	// 9. 建立 Server Snapshot
	if err := app.createServerSnapshot(); err != nil {
		log.Fatal(err)
	}

	// 10. 關聯 Floating IP
	if err := app.associateFloatingIP(networkName); err != nil {
		log.Fatal(err)
	}

	// 11. 刪除所有資源
	if err := app.deleteAllResources(); err != nil {
		log.Fatal(err)
	}
}

func initClient(protocol, baseURL, token, projectCode string) (*vps.Client, *vrm.Client, error) {

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

func (a *App) uploadImageToRepository(imageURL string) error {
	log.Println("")
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

	a.tagID = response.Repository.Tags[0].ID
	a.imageRepositoryID = response.Repository.ID

	log.Printf("Image Repository (%s) Uploaded", a.imageRepositoryID)

	return nil
}

func (a *App) getNetworkAndFlavor() error {
	log.Println("")
	// 2-1: Get default Network ID
	networksList, err := a.vpsClient.Networks().List(a.ctx, &networks.ListNetworksOptions{
		Name: "default",
	})
	if err != nil {
		return err
	}

	if len(networksList) != 1 {
		return errors.New("default network not found")
	}

	a.defaultNetworkID = networksList[0].ID
	log.Printf("Network (%s) Found", a.defaultNetworkID)

	// 2-2: Get first flavor ID from flavor list
	flavorsList, err := a.vpsClient.Flavors().List(a.ctx, &flavors.ListFlavorsOptions{})
	if err != nil {
		return err
	}

	if len(flavorsList) == 0 {
		return errors.New("no flavors found")
	}

	a.firstFlavorID = flavorsList[0].ID
	log.Printf("Flavor (%s) Found", a.firstFlavorID)

	return nil
}

func (a *App) createOrCheckSecurityGroup(securityGroupName string) error {
	log.Println("")
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
		log.Printf("Security group (%s) Created", createdSG.ID)
		a.securityGroupID = createdSG.ID
	}

	return nil
}

func (a *App) createOrCheckKeypair(keypairName string) error {
	log.Println("")

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
		log.Printf("Keypair (%s) Created", a.keypair.ID)
	}

	return nil
}

func (a *App) createServer(serverName, passwordEnvVar string) error {
	log.Println("")
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

	log.Printf("Server (%s) Creating", a.server.ID)

	// Step 7: Wait for server to become active
	// log.Printf("Waiting for server to become active...")
	err = vps.WaitForServerActive(a.ctx, a.vpsClient.Servers(), a.server.ID)
	if err != nil {
		return err
	}
	log.Printf("Server (%s) Actived", a.server.ID)

	return nil
}

func (a *App) createVolume(volumeName string) error {
	log.Println("")
	// Get first volume type
	volumeTypes, err := a.vpsClient.VolumeTypes().List(a.ctx)
	if err != nil {
		return err
	}

	if len(volumeTypes) == 0 {
		return errors.New("no volume types available")
	}

	firstVolumeType := volumeTypes[0]
	log.Printf("Volume type (%s) Found", firstVolumeType)

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

	log.Printf("Non System Volume (%s) Creating", a.volume.ID)

	// Wait for volume to become available
	err = vps.WaitForVolumeAvailable(a.ctx, a.vpsClient.Volumes(), a.volume.ID)
	if err != nil {
		return err
	}
	log.Printf("Non System Volume (%s) Available", a.volume.ID)

	// 7-1: Attach volume to server
	actionReq := &volumes.VolumeActionRequest{
		Action:   volumes.VolumeActionAttach,
		ServerID: a.server.ID,
	}

	err = a.vpsClient.Volumes().Action(a.ctx, a.volume.ID, actionReq)
	if err != nil {
		return err
	}

	log.Printf("Non System Volume (%s) Attaching", a.volume.ID)

	// Wait for volume to become in-use
	err = vps.WaitForVolumeInUse(a.ctx, a.vpsClient.Volumes(), a.volume.ID)
	if err != nil {
		return err
	}
	log.Printf("Non System Volume (%s) In-Use", a.volume.ID)

	return nil
}

func (a *App) createVolumeSnapshot() error {
	log.Println("")
	snapshotName := a.volume.Name + "-snapshot"
	snapshotReq := &snapshots.CreateSnapshotRequest{
		Name:     snapshotName,
		VolumeID: a.volume.ID,
	}

	snapshot, err := a.vpsClient.Snapshots().Create(a.ctx, snapshotReq)
	if err != nil {
		return err
	}

	log.Printf("Volume snapshot (%s) Creating", snapshot.ID)

	// Wait for snapshot to become available
	err = vps.WaitForSnapshotAvailable(a.ctx, a.vpsClient.Snapshots(), snapshot.ID)
	if err != nil {
		return err
	}
	log.Printf("Volume snapshot (%s) Available", snapshot.ID)

	return nil
}

func (a *App) createServerSnapshot() error {
	log.Println("")
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
	log.Printf("Server snapshot repository (%s) Available", a.repositoryID)

	return nil
}

func (a *App) associateFloatingIP(networkName string) error {
	log.Println("")
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

	log.Printf("Floating IP (%s) Associating", fipInfo.ID)

	err = vps.WaitForFloatingIPActive(a.ctx, a.vpsClient.FloatingIPs(), fipInfo.ID)
	if err != nil {
		return err
	}
	log.Printf("Floating IP (%s) Active", fipInfo.ID)

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

	return nil
}

func (a *App) deleteAllResources() error {
	log.Println("")
	log.Println("################################")
	log.Println("#  START DELETE ALL RESOURCES  #")
	log.Println("################################")

	// 11-1: Update Floating IP to Reserved
	log.Println("")
	if a.floatingIP != nil {
		reserved := true
		updateReq := &floatingips.FloatingIPUpdateRequest{
			Reserved: &reserved,
		}

		updatedFIP, err := a.vpsClient.FloatingIPs().Update(a.ctx, a.floatingIP.ID, updateReq)
		if err != nil {
			return err
		}

		log.Printf("Floating IP (%s) Reserved", updatedFIP.ID)
	}

	// 11-2: Disassociate Floating IP
	if a.floatingIP != nil {
		err := a.vpsClient.FloatingIPs().Disassociate(a.ctx, a.floatingIP.ID)
		if err != nil {
			return err
		}

		log.Printf("Floating IP (%s) Disassociated", a.floatingIP.ID)
	}

	// 11-3: Delete Floating IP
	if a.floatingIP != nil {
		err := a.vpsClient.FloatingIPs().Delete(a.ctx, a.floatingIP.ID)
		if err != nil {
			return err
		}
		log.Printf("Floating IP (%s) Deleted", a.floatingIP.ID)
	}

	// 11-4: Delete Volume Snapshots
	log.Println("")
	if a.volume != nil {
		// List all snapshots
		snapshotsList, err := a.vpsClient.Snapshots().List(a.ctx, nil)
		if err != nil {
			return err
		}

		// Delete snapshots that belong to this volume
		for _, snap := range snapshotsList {
			if snap.VolumeID == a.volume.ID {
				err := a.vpsClient.Snapshots().Delete(a.ctx, snap.ID)

				if err != nil {
					log.Printf("Warning: failed to delete snapshot %s: %v", snap.ID, err)
					// Continue with other snapshots, don't fail the entire teardown
				}

				// Wait for snapshot to be neither available nor deleting
				log.Printf("Volume Snapshot (%s) Deleting", snap.ID)
				for {
					// Try to get snapshot status
					_, err := a.vpsClient.Snapshots().Get(a.ctx, snap.ID)
					if err != nil {
						// If snapshot is not found (404), it's been deleted
						log.Printf("Volume Snapshot (%s) Deleted", snap.ID)
						break
					}

					time.Sleep(2 * time.Second)
				}

			}
		}
	}

	// 11-5: Detach non-system Volume
	log.Println("")
	if a.server != nil {
		volumesList, err := a.server.Volumes().List(a.ctx)
		if err != nil {
			return err
		}

		for _, vol := range volumesList {
			if !vol.System {
				err = a.server.Volumes().Detach(a.ctx, vol.VolumeID)
				if err != nil {
					return err
				}
				log.Printf("Non system volume (%s) Detaching", vol.VolumeID)

				// Wait for volume to become available after detach
				err = vps.WaitForVolumeAvailable(a.ctx, a.vpsClient.Volumes(), vol.VolumeID)
				if err != nil {
					return err
				}
				log.Printf("Non system volume (%s) Detached", vol.VolumeID)
			}
		}
	}

	// 11-6: Delete Volume
	if a.volume != nil {
		err := a.vpsClient.Volumes().Delete(a.ctx, a.volume.ID)
		if err != nil {
			return err
		}

		// Wait for volume to be deleted
		log.Printf("Non system Volume (%s) Deleting", a.volume.ID)
		for {
			_, err := a.vpsClient.Volumes().Get(a.ctx, a.volume.ID)
			if err != nil {
				// If volume is not found (404), it's been deleted
				log.Printf("Non system Volume (%s) Deleted", a.volume.ID)
				break
			}
			time.Sleep(2 * time.Second)
		}
	}

	// 11-7: Delete Server
	log.Println("")
	if a.server != nil {
		err := a.vpsClient.Servers().Delete(a.ctx, a.server.ID)
		if err != nil {
			return err
		}
		log.Printf("Server (%s) Deleted", a.server.ID)
	}

	// 11-8: Delete Keypair
	log.Println("")
	if a.keypair != nil {
		err := a.vpsClient.Keypairs().Delete(a.ctx, a.keypair.ID)
		if err != nil {
			return err
		}
		log.Printf("Keypair (%s) Deleted", a.keypair.ID)
	}

	// 11-9: Delete Security Group
	log.Println("")
	if a.securityGroupID != "" {
		err := a.vpsClient.SecurityGroups().Delete(a.ctx, a.securityGroupID)
		if err != nil {
			return err
		}
		log.Printf("Security group (%s) Deleted", a.securityGroupID)
	}

	// 11-10: Delete Repositories
	log.Println("")
	for _, repo := range []string{a.repositoryID, a.imageRepositoryID} {
		if repo != "" {
			err := a.vrmClient.Repositories().Delete(a.ctx, repo)
			if err != nil {
				return err
			}
			log.Printf("Repository (%s) Deleted", repo)
		}
	}

	return nil
}
