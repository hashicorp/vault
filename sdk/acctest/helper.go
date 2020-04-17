package acctest

import (
	"github.com/hashicorp/vault/api"
)

// Helper is intended as a per-package singleton created in TestMain which
// other tests in a package can use to create Terraform execution contexts
type Helper struct {
	// api client for use
	Client *api.Client
	// name for plugin in test
	Name string
	// sourceDir is the dir containing the plugin test binary
	SourceDir string
}

// UseDocker setups docker, copying the plugin test binary
func UseDocker(name, src string) *Helper {
	return &Helper{
		Name:      name,
		SourceDir: src,
	}
}

// Setup creates any temp dir and compiles the binary for copying to Docker
func Setup(name string) error {

	// create temp dir
	// tempDir, err := ioutil.TempDir("", "vault-test-cluster-build")
	// if err != nil {
	// 	return err
	// }

	// wd, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }
	// wd = path.Join(wd, "vault/plugins/uuid")
	// cmd := exec.Command("go", "build", "-o", "./vault/plugins/uuid", "/Users/clint/go-src/github.com/catsby/vault-plugin-secrets-uuid/cmd/uuid/main.go")
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	// err = cmd.Run()
	// if err != nil {
	// 	panic(err)
	// }
	return nil
}
