// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/hcp"
	"github.com/spf13/cobra"
)

var waitForHCPImageReq = &hcp.WaitForImageReq{
	Req: &hcp.GetLatestProductVersionReq{},
}

func newHCPWaitForImageCmd() *cobra.Command {
	availability := ""
	var timeout time.Duration

	waitHCPImage := &cobra.Command{
		Use:   "image",
		Short: "Wait for an HCP image",
		Long:  "Wait for an HCP image",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			waitForHCPImageReq.Req.Availability = hcp.GetLatestProductVersionAvailability(availability)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true // Don't spam the usage on failure

			ctx, cancelCause := context.WithCancelCause(context.Background())
			ctx, cancel := context.WithTimeoutCause(ctx, timeout, errors.New("timed out waiting for image"))
			defer cancel()

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				select {
				case <-ctx.Done():
					return
				case s := <-c:
					fmt.Printf("\x1b[1;33;49mWARNING\x1b[0m: received %s signal. Stopping now..\n", s)
					cancelCause(fmt.Errorf("received signal %s", s))
					cancel()
				}
			}()

			res, err := waitForHCPImageReq.Run(ctx, hcpCmdState.client)
			if err != nil {
				return fmt.Errorf("waiting for an HCP image: %w", err)
			}

			switch rootCfg.format {
			case "json":
				b, err := res.Res.ToJSON()
				if err != nil {
					return err
				}
				fmt.Println(string(b))
			case "markdown":
				tbl := res.Res.ToTable()
				tbl.SetTitle("HCP Image")
				fmt.Println(tbl.RenderMarkdown())
			default:
				fmt.Println(res.Res.ToTable().Render())
			}

			if waitForHCPImageReq.WriteToGithubOutput {
				output, err := res.ToGithubOutput()
				if err != nil {
					return err
				}

				return writeToGithubOutput("wait-image", output)
			}

			return nil
		},
	}

	waitHCPImage.PersistentFlags().StringVarP(&waitForHCPImageReq.Req.ProductName, "product-name", "p", "vault", "The product or component of the image")
	waitHCPImage.PersistentFlags().StringVarP(&waitForHCPImageReq.Req.ProductVersionConstraint, "product-version-constraint", "v", "", "A comma seperated list of constraints. If left unset the latest will be returned")
	waitHCPImage.PersistentFlags().StringVarP(&waitForHCPImageReq.Req.HostManagerVersionConstraint, "host-manager-version-constraint", "m", "", "A semver string. If left unset the latest will be used")
	waitHCPImage.PersistentFlags().StringVarP(&waitForHCPImageReq.Req.CloudProvider, "cloud", "c", "aws", "The cloud provider you wish to search. E.g. aws, azure")
	waitHCPImage.PersistentFlags().StringVarP(&waitForHCPImageReq.Req.CloudRegion, "region", "r", "us-west-2", "The cloud region you wish to search")
	waitHCPImage.PersistentFlags().StringVarP(&availability, "availability", "a", "public", "The image availability")
	waitHCPImage.PersistentFlags().BoolVarP(&waitForHCPImageReq.Req.ExcludeReleaseCandidates, "exclude-release-candidates", "x", false, "Exclude release candidates")
	waitHCPImage.PersistentFlags().DurationVarP(&waitForHCPImageReq.Delay, "delay", "d", 10*time.Second, "the time to wait in-between requests")
	waitHCPImage.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 30*time.Minute, "the maximum duration to wait for the image")
	waitHCPImage.PersistentFlags().BoolVar(&waitForHCPImageReq.WriteToGithubOutput, "github-output", false, "Whether or not to write 'wait-image' to $GITHUB_OUTPUT")

	return waitHCPImage
}
