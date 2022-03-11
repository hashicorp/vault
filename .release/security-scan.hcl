container {
	dependencies = true
	alpine_secdb = true
	secrets      = true
}

binary {
	secrets      = false
	go_modules   = false
	osv          = true
	oss_index    = true
	nvd          = false
}
