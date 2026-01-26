// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var findWorkflowArtifact = &github.FindWorkflowArtifactReq{}

func newGithubFindWorkflowArtifactCmd() *cobra.Command {
	findWorkflowArtifactCmd := &cobra.Command{
		Use:   "workflow-artifact [--pr 1234 | --branch main] [--workflow build --pattern 'vault_[0-9]' ]",
		Short: "Find an artifact associated with a pull requests workflow run",
		Long:  "Find an artifact associated with a pull requests workflow run",
		RunE:  runFindGithubWorkflowArtifactCmd,
	}

	findWorkflowArtifactCmd.PersistentFlags().StringVarP(&findWorkflowArtifact.ArtifactName, "name", "n", "", "The exact artifact name to match")
	findWorkflowArtifactCmd.PersistentFlags().StringVarP(&findWorkflowArtifact.ArtifactPattern, "pattern", "m", "", "A pattern to match an artifact. Only the first match will be returned")
	findWorkflowArtifactCmd.PersistentFlags().StringVarP(&findWorkflowArtifact.Owner, "owner", "o", "hashicorp", "The Github organization")
	findWorkflowArtifactCmd.PersistentFlags().StringVarP(&findWorkflowArtifact.Repo, "repo", "r", "vault", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	findWorkflowArtifactCmd.PersistentFlags().IntVarP(&findWorkflowArtifact.PullNumber, "pr", "p", 0, "The pull request to use as the trigger of the workflow")
	findWorkflowArtifactCmd.PersistentFlags().StringVarP(&findWorkflowArtifact.Branch, "branch", "b", "", "The branch to use as the trigger of the workflow")
	findWorkflowArtifactCmd.PersistentFlags().StringVarP(&findWorkflowArtifact.WorkflowName, "workflow", "w", "", "The name of the workflow the artifact will be associated with")
	findWorkflowArtifactCmd.PersistentFlags().BoolVar(&findWorkflowArtifact.WriteToGithubOutput, "github-output", false, "Whether or not to write 'workflow-artifact' to $GITHUB_OUTPUT")

	return findWorkflowArtifactCmd
}

func runFindGithubWorkflowArtifactCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := findWorkflowArtifact.Run(context.TODO(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("finding workflow artifact: %w", err)
	}

	switch rootCfg.format {
	case "json":
		jsonBytes, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(jsonBytes))
	default:
		fmt.Println(res.ToTable())
	}

	if findWorkflowArtifact.WriteToGithubOutput {
		jsonBytes, err := res.ToGithubOutput()
		if err != nil {
			return err
		}

		return writeToGithubOutput("workflow-artifact", jsonBytes)
	}

	return nil
}
