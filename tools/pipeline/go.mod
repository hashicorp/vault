module github.com/hashicorp/vault/tools/pipeline

go 1.25.11

// We have test modules in here but they ought to be completely ignored
ignore internal/pkg/golang/fixtures

require (
	github.com/Masterminds/semver v1.5.0
	github.com/PuerkitoBio/goquery v1.12.0
	github.com/avast/retry-go/v4 v4.7.0
	github.com/google/go-github/v83 v83.0.0
	github.com/hashicorp/hcl/v2 v2.24.0
	github.com/hashicorp/releases-api v0.4.3
	github.com/jedib0t/go-pretty/v6 v6.8.2
	github.com/owenrumney/go-sarif/v3 v3.3.0
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2
	github.com/shurcooL/githubv4 v0.0.0-20260209031235-2402fdf4a9ed
	github.com/slack-go/slack v0.27.0
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.11.1
	github.com/veqryn/slog-context v0.9.0
	github.com/zclconf/go-cty v1.19.0
	golang.org/x/mod v0.37.0
	golang.org/x/oauth2 v0.36.0
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/andybalholm/cascadia v1.3.4 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/apparentlymart/go-textseg/v17 v17.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.25.3 // indirect
	github.com/go-openapi/errors v0.22.8 // indirect
	github.com/go-openapi/jsonpointer v0.24.0 // indirect
	github.com/go-openapi/jsonreference v0.21.6 // indirect
	github.com/go-openapi/loads v0.24.0 // indirect
	github.com/go-openapi/runtime v0.32.4 // indirect
	github.com/go-openapi/runtime/server-middleware v0.32.4 // indirect
	github.com/go-openapi/spec v0.22.6 // indirect
	github.com/go-openapi/strfmt v0.26.4 // indirect
	github.com/go-openapi/swag v0.27.0 // indirect
	github.com/go-openapi/swag/cmdutils v0.27.0 // indirect
	github.com/go-openapi/swag/conv v0.27.0 // indirect
	github.com/go-openapi/swag/fileutils v0.27.0 // indirect
	github.com/go-openapi/swag/jsonname v0.27.0 // indirect
	github.com/go-openapi/swag/jsonutils v0.27.0 // indirect
	github.com/go-openapi/swag/loading v0.27.0 // indirect
	github.com/go-openapi/swag/mangling v0.27.0 // indirect
	github.com/go-openapi/swag/netutils v0.27.0 // indirect
	github.com/go-openapi/swag/stringutils v0.27.0 // indirect
	github.com/go-openapi/swag/typeutils v0.27.0 // indirect
	github.com/go-openapi/swag/yamlutils v0.27.0 // indirect
	github.com/go-openapi/validate v0.26.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jessevdk/go-flags v1.6.1 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/mattn/go-runewidth v0.0.24 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/mitchellh/pointerstructure v1.2.1 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/shurcooL/graphql v0.0.0-20240915155400-7ee5256398cf // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.39.0 // indirect
	golang.org/x/tools v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
