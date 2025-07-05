// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
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
		Short: "Create a list of Vault versions between a lower and upper bound",
		Long:  "Create a list of Vault versions between a lower and upper bound",
		RunE:  runListVersionsReq,
	}

	versionsCmd.PersistentFlags().StringVarP(&listReleaseVersionsReq.UpperBound, "upper", "u", "", "The highest version to include")
	versionsCmd.PersistentFlags().StringVarP(&listReleaseVersionsReq.LowerBound, "lower", "l", "", "The lowest version to include")
	versionsCmd.PersistentFlags().UintVarP(&listReleaseVersionsReq.NMinus, "nminus", "n", 0, "Instead of setting a dedicated lower bound, calculate N-X from the upper")
	versionsCmd.PersistentFlags().StringVarP(&listReleaseVersionsReq.LicenseClass, "edition", "e", "", "The edition of Vault. Can either be 'ce' or 'enterprise'")
	versionsCmd.PersistentFlags().StringSliceVarP(&listReleaseVersionsReq.Skip, "skip", "s", []string{}, "Skip this version. Can be provided none-to-many times")

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

	res, err := listReleaseVersionsReq.Run(context.TODO())
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
