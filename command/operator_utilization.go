// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/mapstructure"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorUtilizationCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorUtilizationCommand)(nil)
)

type OperatorUtilizationCommand struct {
	*BaseCommand

	flagMessage   string
	flagTodayOnly BoolPtr
	flagOutput    string
}

func (c *OperatorUtilizationCommand) Synopsis() string {
	return "Generates license utilization reporting bundle"
}

func (c *OperatorUtilizationCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *OperatorUtilizationCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorUtilizationCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "message",
		Target:     &c.flagMessage,
		Completion: complete.PredictAnything,
		Usage:      "Provide context about the conditions under which the report was generated and submitted. This message is not included in the license utilization bundle but will be included in the vault server logs.",
	})

	f.BoolPtrVar(&BoolPtrVar{
		Name:   "today-only",
		Target: &c.flagTodayOnly,
		Usage:  "To include only today’s snapshot, no historical snapshots. If no snapshots were persisted in the last 24 hrs, it takes a snapshot and exports it to a bundle.",
	})

	f.StringVar(&StringVar{
		Name:       "output",
		Target:     &c.flagOutput,
		Completion: complete.PredictAnything,
		Usage:      "Specifies the output path for the bundle. Defaults to a time-based generated file name.",
	})

	return set
}

func (c *OperatorUtilizationCommand) Help() string {
	helpText := `
Usage: vault operator utilization [options]

Produces a bundle of snapshots that contains license utilization data. If no snapshots were persisted in the last 24 hrs, it takes a snapshot and includes it in the bundle to prevent stale data.

  To create a license utilization bundle that includes all persisted historical snapshots and has the default bundle name:
  
  $ vault operator utilization

  To create a license utilization bundle with a message about the bundle (Note: this message is not included in the bundle but only included in server logs): 

  $ vault operator utilization -message="Change Control 654987"

  To create a license utilization bundle with only today's snapshot:

  $ vault operator utilization -today-only

  To create a license utilization bundle with a specific name:

  $ vault operator utilization -output="/utilization/reports/latest.json"

` + c.Flags().Help()

	return helpText
}

func (c *OperatorUtilizationCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	parsedArgs := f.Args()
	if len(parsedArgs) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(parsedArgs)))
		return 1
	}

	outputBundleFile, err := getOutputFileName(time.Now().UTC(), c.flagOutput)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error during validation: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Capture license utilization reporting data
	bundleDataBytes, err := c.getManualReportingCensusData(client)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error capturing license utilization reporting data: %s", err))
		return 1
	}

	err = os.WriteFile(outputBundleFile, bundleDataBytes, 0o400)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing license utilization reporting data to bundle %q: %s", outputBundleFile, err))
		return 1
	}

	c.UI.Info(fmt.Sprintf("Success! License utilization reporting bundle written to: %s", outputBundleFile))
	return 0
}

// getOutputFileName returns the file name of the license utilization reporting bundle ending with .json
// If filename is a path with non-existing parent directory, it creates a new directory to which the file with returned filename is added
func getOutputFileName(inputTime time.Time, flagOutput string) (string, error) {
	formattedTime := inputTime.Format(fileFriendlyTimeFormat)
	switch len(flagOutput) {
	case 0:
		flagOutput = fmt.Sprintf("vault-utilization-%s.json", formattedTime)
	default:
		flagOutput = filepath.Clean(flagOutput)
		ext := filepath.Ext(flagOutput)
		switch ext {
		case "": // it's a directory
			flagOutput = filepath.Join(flagOutput, fmt.Sprintf("vault-utilization-%s.json", formattedTime))
		case ".json":
		default:
			return "", fmt.Errorf("invalid file extension %s, must be .json", ext)
		}
	}

	// Stat the file to ensure we don't override any existing data.
	_, err := os.Stat(flagOutput)
	switch {
	case os.IsNotExist(err):
	case err != nil:
		return "", fmt.Errorf("unable to stat file: %s", err)
	default:
		return "", fmt.Errorf("output file already exists: %s", flagOutput)
	}

	// output file does not exist, create the parent directory if it doesn't exist
	_, err = os.Stat(filepath.Dir(flagOutput))
	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(filepath.Dir(flagOutput), 0o700)
		if err != nil {
			return "", fmt.Errorf("unable to create output directory: %s", err)
		}
	case err != nil:
		return "", fmt.Errorf("unable to stat directory: %s", err)
	}
	return flagOutput, nil
}

func (c *OperatorUtilizationCommand) getManualReportingCensusData(client *api.Client) ([]byte, error) {
	data := make(map[string]interface{})
	if c.flagTodayOnly.IsSet() {
		data["today_only"] = c.flagTodayOnly.Get()
	}
	if c.flagMessage != "" {
		data["message"] = c.flagMessage
	}
	secret, err := client.Logical().Write("sys/utilization", data)
	if err != nil {
		return nil, fmt.Errorf("error getting license utilization reporting data: %w", err)
	}
	if secret == nil {
		return nil, errors.New("no license utilization reporting data available")
	}

	var bundleBase64Str string
	err = mapstructure.Decode(secret.Data["utilization_bundle"], &bundleBase64Str)
	if err != nil {
		return nil, err
	}

	bundleByteArray, err := base64.StdEncoding.DecodeString(bundleBase64Str)
	if err != nil {
		return nil, err
	}
	return bundleByteArray, nil
}
