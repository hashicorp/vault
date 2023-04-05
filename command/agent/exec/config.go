package exec

type EnvTemplate struct {
	Name            string `hcl:"name,label"`
	Contents        string `hcl:"contents"`
	ErrOnMissingKey bool   `hcl:"err_on_missing_key"`
	Group           string `hcl:"group"`
}
