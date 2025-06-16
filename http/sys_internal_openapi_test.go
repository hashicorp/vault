// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func loadExpectedSpec(t *testing.T, path string) map[string]interface{} {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open expected spec file: %v", err)
	}
	defer file.Close()

	var spec map[string]interface{}
	if err := json.NewDecoder(file).Decode(&spec); err != nil {
		t.Fatalf("failed to decode expected spec: %v", err)
	}

	return spec
}

func genOpenAPISpec(t *testing.T) map[string]interface{} {
	// Path to the gen_openapi.sh script
	scriptPath := "../scripts/gen_openapi.sh"

	// Create the command to execute the script
	cmd := exec.Command("bash", scriptPath)

	// Capture the output and error streams
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the script
	err := cmd.Run()
	if err != nil {
		t.Fatalf("failed to run gen_openapi.sh: %v\nstderr: %s", err, stderr.String())
	}

	// Optionally, check the output
	t.Logf("gen_openapi.sh output:\n%s", stdout.String())

	return loadExpectedSpec(t, "../scripts/openapi.json")
}

func Test_OpenAPISpecResponse(t *testing.T) {
	// Path to the expected OpenAPI spec file
	expectedSpecPath := "../expected.json"

	// Load the expected spec
	expectedSpec := loadExpectedSpec(t, expectedSpecPath)

	// Generate the OpenAPI spec (replace `GenerateOpenAPISpec` with the actual function)
	generatedSpecMap := genOpenAPISpec(t)

	// Decode the generated spec into a comparable structure
	// var generatedSpecMap map[string]interface{}
	// if err := json.Unmarshal([]byte(generatedSpec), &generatedSpecMap); err != nil {
	// 	t.Fatalf("failed to decode generated spec: %v", err)
	// }

	// Compare the specs
	if diff := deep.Equal(expectedSpec, generatedSpecMap); diff != nil {
		t.Errorf("OpenAPI spec mismatch:\n%s", strings.Join(diff, "\n"))
	}
}
