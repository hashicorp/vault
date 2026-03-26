// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var listReleaseVersionsReq = &releases.ListVersionsReq{
	VersionLister: releases.NewClient(),
}

func newReleasesVersionsBetweenCmd() *cobra.Command {
	versionsCmd := &cobra.Command{
		Use:   "versions",
		Short: "Create a list of product versions between a lower and upper bound",
		Long:  "Create a list of product versions between a lower and upper bound",
		RunE:  runListVersionsReq,
	}

	versionsCmd.PersistentFlags().StringVarP(&listReleaseVersionsReq.UpperBound, "upper", "u", "", "The highest version to include")
	versionsCmd.PersistentFlags().StringVarP(&listReleaseVersionsReq.LowerBound, "lower", "l", "", "The lowest version to include")
	versionsCmd.PersistentFlags().UintVarP(&listReleaseVersionsReq.NMinus, "nminus", "n", 0, "Instead of setting a dedicated lower bound, calculate N-X from the upper")
	versionsCmd.PersistentFlags().StringVarP(&listReleaseVersionsReq.LicenseClass, "edition", "e", "", "The edition of the product. Can either be 'ce' or 'enterprise'")
	versionsCmd.PersistentFlags().StringSliceVarP(&listReleaseVersionsReq.Skip, "skip", "s", []string{}, "Skip this version. Can be provided none-to-many times")

	// Cadence flags for version calculation
	versionsCmd.PersistentFlags().StringVar((*string)(&listReleaseVersionsReq.Cadence), "cadence", "major", "Version cadence for n-minus calculation: 'minor' or 'major' (default: major)")
	versionsCmd.PersistentFlags().StringVar(&listReleaseVersionsReq.TransitionVersion, "transition-version", "", "Last version of previous cadence (optional, for cadence transitions)")
	versionsCmd.PersistentFlags().StringVar((*string)(&listReleaseVersionsReq.PriorCadence), "prior-cadence", "", "If the product transitioned from an older release candence to a new one, this defines the prior cadence type: 'minor' or 'major' (required if transition-version is set)")
	versionsCmd.PersistentFlags().StringVar(&listReleaseVersionsReq.ProductName, "product", "vault", "Product name")

	err := versionsCmd.MarkPersistentFlagRequired("upper")
	if err != nil {
		panic(err)
	}
	err = versionsCmd.MarkPersistentFlagRequired("edition")
	if err != nil {
		panic(err)
	}

	return versionsCmd
}

func runListVersionsReq(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := listReleaseVersionsReq.Run(cmd.Context())
	if err != nil {
		return err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return err
}
