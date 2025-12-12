module github.com/hashicorp/vault/pipeline/golang/modb

go 1.25.2

godebug (
	default=go1.21
	httpcookiemaxnum=4000
	panicnil=1
)

exclude golang.org/x/term v0.2.0

replace github.com/99designs/keyring => github.com/Jeffail/keyring v1.2.3

require github.com/99designs/keyring v0.0.0-00010101000000-000000000000

retract (
	[v1.0.0, v1.9.9]
	v0.9.0
)

tool golang.org/x/tools/cmd/stringer

require (
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
)

ignore ./node_modules
