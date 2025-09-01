package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

func main() {
	os.Setenv("OCI_CONFIG_FILE", "./config")

	configProvider := common.DefaultConfigProvider()

	client, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		panic(fmt.Sprintf("Failed to create compute client: %v", err))
	}

	compartmentID := os.Getenv("COMPARTMENT_ID")
	availabilityDomain := os.Getenv("AVAILABILITY_DOMAIN")
	subnetID := os.Getenv("SUBNET_ID")
	imageID := os.Getenv("IMAGE_ID")
	sshKey := os.Getenv("SSH_KEY")
	name := os.Getenv("INSTANCE_NAME")
	shape := os.Getenv("INSTANCE_SHAPE")

	instanceDetails := core.LaunchInstanceDetails{
		CompartmentId:      &compartmentID,
		AvailabilityDomain: &availabilityDomain,
		Shape:              &shape,
		DisplayName:        &name,
		CreateVnicDetails: &core.CreateVnicDetails{
			SubnetId:       &subnetID,
			AssignPublicIp: common.Bool(true),
		},
		Metadata: map[string]string{
			"ssh_authorized_keys": sshKey,
		},
		SourceDetails: core.InstanceSourceViaImageDetails{
			ImageId: &imageID,
		},
		ShapeConfig: &core.LaunchInstanceShapeConfigDetails{
			Ocpus:       common.Float32(4.0),
			MemoryInGBs: common.Float32(24.0),
		},
	}

	ctx := context.Background()
	for {
		resp, err := client.LaunchInstance(ctx, core.LaunchInstanceRequest{
			LaunchInstanceDetails: instanceDetails,
		})
		if err != nil {
			if serviceErr, ok := err.(common.ServiceError); ok {
				fmt.Printf("Got OCI Service Error: %s\n", serviceErr.GetMessage())
				fmt.Printf("Status Code: %d\n", serviceErr.GetHTTPStatusCode())

				if serviceErr.GetHTTPStatusCode() == 429 {
					fmt.Println("❌ Too many requests: sleeping...")
					time.Sleep(3 * time.Minute)
				} else {
					fmt.Printf("⚠️ HTTP Status Code: %d\n", serviceErr.GetHTTPStatusCode())
					time.Sleep(1 * time.Minute)
				}
			} else {
				// Not an OCI service error (e.g., network issue, timeout)
				fmt.Printf("Non-OCI error: %v\n", err)
			}
			continue
		}

		fmt.Printf("Launched instance with OCID: %s\n", *resp.Instance.Id)
		return
	}
}
