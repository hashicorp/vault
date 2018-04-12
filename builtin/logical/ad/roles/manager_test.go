package roles

import "testing"

func TestPathRegexpWorks(t *testing.T) {
	// TODO write a test based on this:
	/*
		func main() {
			tests := []string{
				"roles",
				"roles/",
				"rolessuper",
				"roles/candy",
				"cats/roles",
				"roles/beccas_role",
			}

			pattern := `^roles$|^roles/.|^roles/$`

			re, err := regexp.Compile(pattern)
			if err != nil {
				fmt.Println("err: " + err.Error())
				return
			}

			for _, test := range tests {
				matches := re.FindStringSubmatch(test)
				if matches == nil {
					fmt.Println("discarding " + test)
					continue
				}
				fmt.Println(test + " is considered a match")
			}
		}
	*/
	// also, test it starts with ^ and ends with $ because if it doesn't,
	// the outer framework will add it on and derp your beautiful regexp
}
