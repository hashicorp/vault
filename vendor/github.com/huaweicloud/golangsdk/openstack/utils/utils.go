package utils

func DeleteNotPassParams(params *map[string]interface{}, not_pass_params []string) {
	for _, i := range not_pass_params {
		delete(*params, i)
	}
}
