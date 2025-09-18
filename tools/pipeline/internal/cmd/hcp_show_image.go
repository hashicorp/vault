// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/hcp"
	"github.com/spf13/cobra"
)

var showHCPImageReq = &hcp.GetLatestProductVersionReq{}

func newHCPShowImageCmd() *cobra.Command {
	availability := ""

	showHCPImage := &cobra.Command{
		Use:   "image",
		Short: "Show details of an HCP image",
		Long:  "Show details of an HCP image",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			showHCPImageReq.Availability = hcp.GetLatestProductVersionAvailability(availability)
		},
		RunE: runHCPImageShowLatestCmd,
	}

	showHCPImage.PersistentFlags().StringVarP(&showHCPImageReq.ProductName, "product-name", "p", "vault", "The product or component of the image")
	showHCPImage.PersistentFlags().StringVarP(&showHCPImageReq.ProductVersionConstraint, "product-version-constraint", "v", "", "A comma seperated list of constraints. If left unset the latest will be returned")
	showHCPImage.PersistentFlags().StringVarP(&showHCPImageReq.HostManagerVersionConstraint, "host-manager-version-constraint", "m", "", "A semver string. If left unset the latest will be used")
	showHCPImage.PersistentFlags().StringVarP(&showHCPImageReq.CloudProvider, "cloud", "c", "aws", "The cloud provider you wish to search. E.g. aws, azure")
	showHCPImage.PersistentFlags().StringVarP(&showHCPImageReq.CloudRegion, "region", "r", "us-west-2", "The cloud region you wish to search")
	showHCPImage.PersistentFlags().StringVarP(&availability, "availability", "a", "public", "The image availability")
	showHCPImage.PersistentFlags().BoolVarP(&showHCPImageReq.ExcludeReleaseCandidates, "exclude-release-candidates", "x", false, "Exclude release candidates")

	return showHCPImage
}

func runHCPImageShowLatestCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := showHCPImageReq.Run(context.TODO(), hcpCmdState.client)
	if err != nil {
		return fmt.Errorf("showing HCP image: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "markdown":
		tbl := res.ToTable()
		tbl.SetTitle("HCP Image")
		fmt.Println(tbl.RenderMarkdown())
	default:
		fmt.Println(res.ToTable().Render())
	}

	return nil
}
