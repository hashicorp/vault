// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/sarif"
	"github.com/spf13/cobra"
)

var uploadUSTReq = &sarif.UploadUSTArchiveReq{}

func newSarifUploadUSTArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-ust [sarif-file]",
		Short: "Upload a SARIF file to UST",
		Args:  cobra.ExactArgs(1),
		RunE:  runSarifUploadUSTCmd,
		Long: `Convert a SARIF file to Concert format and publish it to the UST
(Unified Security Tracking) destination Git repository via a push.

Authentication uses a pre-minted GitHub token supplied via the GITHUB_TOKEN or
GH_TOKEN environment variable. In a GitHub Actions workflow, use
actions/create-github-app-token to mint a short-lived installation token and
set GH_TOKEN to its output before invoking this command.

The Concert JSON is packaged as a dynamic-scan archive and pushed to the
results/<product-id>/<squad-id>/ path on the destination branch.

Examples:
  # Upload a ZAP SARIF file to UST (token provided via environment)
  GH_TOKEN="$(gh auth token)" pipeline sarif upload-ust report.json \
    --offering-id vault \
    --product-id vault-enterprise \
    --squad-id vault-ent \
    --source-repo-url https://github.com/hashicorp/vault-enterprise \
    --source-branch main \
    --source-commit abcd1234 \
    --org ibm-ust
    --dest-branch main

  # Also emit the Concert JSON to a file
  pipeline sarif upload-ust report.json ... --out concert.json`,
	}

	cmd.Flags().StringVar(&uploadUSTReq.OfferingID, "offering-id", "", "UST offering ID")
	cmd.Flags().StringVar(&uploadUSTReq.ProductID, "product-id", "", "UST product ID (destination repo name)")
	cmd.Flags().StringVar(&uploadUSTReq.SquadID, "squad-id", "", "UST squad ID")
	cmd.Flags().StringVar(&uploadUSTReq.SourceRepoURL, "source-repo-url", "", "URL of the source repository")
	cmd.Flags().StringVar(&uploadUSTReq.SourceBranch, "source-branch", "", "Source repository branch")
	cmd.Flags().StringVar(&uploadUSTReq.SourceCommit, "source-commit", "", "Source repository commit SHA")
	cmd.Flags().StringVar(&uploadUSTReq.Host, "host", "github.ibm.com", "GitHub host for UST destination")
	cmd.Flags().StringVar(&uploadUSTReq.Org, "org", "", "GitHub organization owning the UST destination repository")
	cmd.Flags().StringVar(&uploadUSTReq.DestBranch, "dest-branch", "main", "Destination branch in UST repo")
	cmd.Flags().StringVar(&uploadUSTReq.GitName, "git-name", "hc-github-team-secure-vault-core", "Git committer name")
	cmd.Flags().StringVar(&uploadUSTReq.GitEmail, "git-email", "github-team-secure-vault-core@hashicorp.com", "Git committer email")
	cmd.Flags().StringVarP(&uploadUSTReq.ConcertOutPath, "out", "o", "", "Optional: also write Concert JSON to this path")
	cmd.Flags().IntVar(&uploadUSTReq.PushRetries, "push-retries", 1, "Number of pull-rebase+push retries after a push conflict")

	var err error
	for _, f := range []string{
		"org",
		"offering-id",
		"product-id",
		"squad-id",
		"source-repo-url",
		"source-branch",
		"source-commit",
	} {
		if err = cmd.MarkFlagRequired(f); err != nil {
			panic(err)
		}
	}

	return cmd
}

func runSarifUploadUSTCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	uploadUSTReq.SarifPath = args[0]
	res, err := uploadUSTReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("uploading to UST: %w", err)
	}

	var output []byte
	switch rootCfg.format {
	case "json":
		output, err = res.ToJSON()
		if err != nil {
			return fmt.Errorf("formatting output as JSON: %w", err)
		}
	case "markdown":
		output = []byte(res.ToMarkdown())
	default:
		output = []byte(res.ToTable().Render())
	}

	if _, err := os.Stdout.Write(output); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
