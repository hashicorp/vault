package cmd

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/vault/helper/oas"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/vault"
)

func Run() int {
	doc := oas.NewOASDoc()

	// we can choose to build different things at this point
	buildDoc(&doc)

	// we can choose to render different things at this point
	//oas, err := apidoc.NewoasRenderer(2)

	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	//	return 1
	//}
	//oas.Render(doc)
	output, _ := json.Marshal(doc)
	os.Stdout.Write(output)

	return 0
}

// buildDoc is a sample of how to populate a Document with content
// from backends or other sources.
func buildDoc(doc *oas.OASDoc) {
	// Load the /sys backend, and then append the separate manual paths
	backend := vault.NewSystemBackend(&vault.Core{}, nil).Backend
	framework.DocumentPaths(backend, doc)
	//doc.AddPath("sys", vault.ManualPaths()...)

	// Load another backend to show how separate mounts could be presented.
	// This will be in a separate, "aws" group in the output oas.
	//framework.LoadBackend(aws.Backend().Backend, doc)
}
