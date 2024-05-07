**Table of Contents**

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Vault UI](#vault-ui)
  - [Ember CLI Version Upgrade Matrix](#ember-cli-version-upgrade-matrix)
  - [Prerequisites](#prerequisites)
  - [Running a Vault Server](#running-a-vault-server)
  - [Running the UI locally](#running-the-ui-locally)
    - [Mirage](#mirage)
  - [Building Vault UI into a Vault Binary](#building-vault-ui-into-a-vault-binary)
  - [Development](#development)
    - [Quick commands](#quick-commands)
    - [Code Generators](#code-generators)
    - [Running Tests](#running-tests)
    - [Linting](#linting)
    - [Contributing / Best Practices](#contributing--best-practices)
  - [Further Reading / Useful Links](#further-reading--useful-links)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Vault UI

This README outlines the details of collaborating on this Ember application.

## Ember CLI Version Upgrade Matrix

| Vault Version | Ember Version |
| ------------- | ------------- |
| 1.17.x        | 5.4.2         |
| 1.15.x        | 4.12.0        |
| 1.14.x        | 4.4.0         |
| 1.13.x        | 4.4.0         |
| 1.12.x        | 3.28.5        |
| 1.11.x        | 3.28.5        |
| 1.10.x        | 3.28.5        |
| 1.9.x         | 3.22.0        |
| 1.8.x         | 3.22.0        |
| 1.7.x         | 3.11.0        |

## Prerequisites

You will need the following things properly installed on your computer.

- [Git](https://git-scm.com/)
- [Node.js](https://nodejs.org/)
- [Yarn](https://yarnpkg.com/)
- [Ember CLI](https://cli.emberjs.com/release/)
- [Google Chrome](https://google.com/chrome/)

In order to enforce the same version of `yarn` across installs, the `yarn` binary is included in the repo
in the `.yarn/releases` folder. To update to a different version of `yarn`, use the `yarn policies set-version VERSION` command. For more information on this, see the [documentation](https://yarnpkg.com/en/docs/cli/policies).

## Running a Vault Server

Before running Vault UI locally, a Vault server must be running. First, ensure
Vault dev is built according the instructions in `../README.md`.

- To start a single local Vault server: `yarn vault`
- To start a local Vault cluster: `yarn vault:cluster`

These commands may also be [aliased on your local device](https://github.com/hashicorp/vault-tools/blob/master/users/noelle/vault_aliases).

## Running the UI locally

To spin up the UI, a Vault server must be running (see previous step).
_All of the commands below assume you're in the `ui/` directory._

> These steps will start an Ember CLI server that proxies requests to port 8200,
> and enable live rebuilding of the application as you change the UI application code.
> Visit your app at [http://localhost:4200](http://localhost:4200).

1. Install dependencies:

`yarn`

2. Run Vault UI and proxy back to a Vault server running on the default port, 8200:

`yarn start`

> If your Vault server is running on a different port you can use the
> long-form version of the npm script:

`ember server --proxy=http://localhost:PORT`

### Mirage

[Mirage](https://miragejs.com/docs/getting-started/introduction/) can be helpful for mocking backend endpoints.
Look in [mirage/handlers](mirage/handlers/) for existing mocked backends.

Run yarn with mirage: `export MIRAGE_DEV_HANDLER=<handler> yarn start`

Where `handlername` is one of the options exported in [mirage/handlers/index](mirage/handlers/index.js)

## Building Vault UI into a Vault Binary

We use the [embed](https://golang.org/pkg/embed/) package from Go >1.20 to build
the static assets of the Ember application into a Vault binary.

This can be done by running these commands from the root directory:
`make static-dist`
`make dev-ui`

This will result in a Vault binary that has the UI built-in - though in
a non-dev setup it will still need to be enabled via the `ui` config or
setting `VAULT_UI` environment variable.

## Development

### Quick commands

| Command                                           | Description                                                             |
| ------------------------------------------------- | ----------------------------------------------------------------------- |
| `yarn start`                                      | start the app with live reloading (vault must be running on port :8200) |
| `export MIRAGE_DEV_HANDLER=<handler>; yarn start` | start the app with the mocked mirage backend, with handler provided     |
| `make static-dist && make dev-ui`                 | build a Vault binary with UI assets (run from root directory not `/ui`) |
| `ember g component foo -ir core`                  | generate a component in the /addon engine                               |
| `yarn test:filter`                                | run non-enterprise in the browser                                       |
| `yarn test:filter -f='<test name>'`               | run tests in the browser, filtering by test name                        |
| `yarn lint:js`                                    | lint javascript files                                                   |

### Code Generators

Make use of the many generators for code, try `ember help generate` for more details. If you're using a component that can be widely-used, consider making it an `addon` component instead (see [this PR](https://github.com/hashicorp/vault/pull/6629) for more details)

eg. a reusable component named foo that you'd like in the core engine (read more about Ember engines [here](https://ember-engines.com/docs)).

- `ember g component foo -ir core`

The above command creates a template-only component by default. If you'd like to add a backing class, add the `-gc` flag:

- `ember g component foo -gc -ir core`

### Running Tests

Running tests will spin up a Vault dev server on port :9200 via a pretest script that testem (the test runner) executes. All of the acceptance tests then run, which proxy requests back to that server. The normal test scripts use `ember-exam` which split into parallel runs, which is excellent for speed but makes it harder to debug. So we have a custom yarn script that automatically opens all the tests in a browser, and we can pass the `-f` flag to target the test(s) we're debugging.

- `yarn run test` lint & run all the tests (CI uses this)
- `yarn run test:oss` lint & run all the non-enterprise tests (CI uses this)
- `yarn run test:quick` run all the tests without linting
- `yarn run test:quick-oss` run all the non-enterprise tests without linting
- `yarn run test:filter -f="policies"` run the filtered test in the browser with no splitting. `-f` is set to `!enterprise` by default
  [QUnit's `filter` config](https://api.qunitjs.com/config/QUnit.config#qunitconfigfilter-string--default-undefined)

### Linting

- `yarn lint:js`
- `yarn lint:hbs`
- `yarn lint:fix`

### Contributing / Best Practices

Hello and thank you for contributing to the Vault UI! Below is a list of patterns we follow on the UI team to keep in mind when contributing to the UI codebase. This is an ever-evolving process, so we welcome any comments, questions or general feedback.

> **Remember** prefixing your branch name with `ui/` will run UI tests and skip the go tests. If your PR includes backend changes, _do not_ prefix your branch, instead add the `ui` label on github. This will trigger the UI test suite to run, in addition to the backend Go tests.

- [routing](docs/routing.md)
- [serializers/adapters](docs/serializers-adapters.md)
- [models](docs/models.md)
- [components](docs/components.md)
- [forms](docs/forms.md)
- [css](docs/css.md)
- [ember engines](docs/engines.md)

## Further Reading / Useful Links

- [ember.js](https://emberjs.com/)
- [ember-cli](https://cli.emberjs.com/release/)
- Development Browser Extensions
  - [ember inspector for chrome](https://chrome.google.com/webstore/detail/ember-inspector/bmdblncegkenkacieihfhpjfppoconhi)
  - [ember inspector for firefox](https://addons.mozilla.org/en-US/firefox/addon/ember-inspector/)
