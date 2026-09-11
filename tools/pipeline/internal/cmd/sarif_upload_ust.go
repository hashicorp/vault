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
		Long: `Convert a SARIF file to Concert format and push it to the UST
(Unified Security Tracking) destination Git repository.

The command reads a GitHub token from GITHUB_TOKEN or GH_TOKEN. In a GitHub
Actions workflow, use actions/create-github-app-token to mint a short-lived
installation token and export it via one of those variables before running
this command.

The Concert JSON is packaged as a dynamic-scan archive and committed to
results/<product-id>/<squad-id>/ on the destination branch. Since the upload
step will push to the same branch, if you're running several of these in
parallel you'll likely have jobs win and lose the push race. Use the
--upload-retries flag to gracefully retry the upload.

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
	cmd.Flags().IntVar(&uploadUSTReq.UploadRetries, "upload-retries", 1, "Max retry attempts for the pull-rebase+push loop when jobs race to push")

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
