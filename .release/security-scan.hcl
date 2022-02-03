container {
	dependencies = true
	alpine_secdb = true
	secrets      = true
}

binary {
	secrets      = true
	go_modules   = false
	osv          = true
	oss_index    = true
	nvd          = true
}
