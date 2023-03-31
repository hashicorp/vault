**Table of Contents**

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Vault UI](#vault-ui)
  - [Ember CLI Version Matrix](#ember-cli-version-matrix)
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
    - [Further Reading / Useful Links](#further-reading--useful-links)
  - [Best Practices](#best-practices)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Vault UI

This README outlines the details of collaborating on this Ember application.

## Ember CLI Version Matrix

| Vault Version | Ember Version |
| ------------- | ------------- |
| 1.13.x        | 4.4.0         |
| 1.10.x        | 3.28.5        |
| 1.9.x         | 3.22.0        |
| 1.8.x         | 3.22.0        |
| 1.7.x         | 3.11          |

## Prerequisites

You will need the following things properly installed on your computer.

- [Git](https://git-scm.com/)
- [Node.js](https://nodejs.org/)
- [Yarn](https://yarnpkg.com/)
- [Ember CLI](https://cli.emberjs.com/release/)
- [Google Chrome](https://google.com/chrome/)

* [lint-staged\*](https://www.npmjs.com/package/lint-staged)

\* lint-staged is an optional dependency - running `yarn` will install it.
If don't want optional dependencies installed you can run `yarn --ignore-optional`. If you've ignored the optional deps
previously and want to install them, you have to tell yarn to refetch all deps by running `yarn --force`.

In order to enforce the same version of `yarn` across installs, the `yarn` binary is included in the repo
in the `.yarn/releases` folder. To update to a different version of `yarn`, use the `yarn policies set-version VERSION` command. For more information on this, see the [documentation](https://yarnpkg.com/en/docs/cli/policies).

## Running a Vault Server

Before running Vault UI locally, a Vault server must be running. First, ensure
Vault dev is built according the the instructions in `../README.md`. To start a
single local Vault server:

- `yarn vault`

To start a local Vault cluster:

- `yarn vault:cluster`

These commands may also be [aliased on your local device](https://github.com/hashicorp/vault-tools/blob/master/users/noelle/vault_aliases).

## Running the UI locally

> a Vault server must be running, see step above.
> These steps will start an Ember CLI server that proxies requests to port 8200,
> and enable live rebuilding of the application as you change the UI application code.
> Visit your app at [http://localhost:4200](http://localhost:4200).

To get all of the JavaScript dependencies installed, run this in the `ui` directory:

- `yarn`

If you want to run Vault UI and proxy back to a Vault server running
on the default port, 8200, run the following in the `ui` directory:

- `yarn start`

If your Vault server is running on a different port you can use the
long-form version of the npm script:

`ember server --proxy=http://localhost:PORT`

### Mirage

To run yarn with mirage, do:

- `yarn start:mirage handlername`

Where `handlername` is one of the options exported in `mirage/handlers/index`

### Building Vault UI into a Vault Binary

We use the [embed](https://golang.org/pkg/embed/) package from Go >1.20 to build
the static assets of the Ember application into a Vault binary.

This can be done by running these commands from the root directory:
`make static-dist`
`make dev-ui`

This will result in a Vault binary that has the UI built-in - though in
a non-dev setup it will still need to be enabled via the `ui` config or
setting `VAULT_UI` environment variable.

## Development

## Quick commands

| Command                                                                                  | Description                                                         |
| ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| `yarn start`                                                                             | start the app with live reloading                                   |
| `yarn start:mirage <handler>`                                                            | start the app with the mocked mirage backend, with handler provided |
| `make static-dist && make dev-ui`                                                        | build a Vault binary with UI assets, by default runs on port :8200  |
| `ember g component foo --in lib/core`                                                    | generate a component in the /addon engine                           |
| `echo "export { default } from 'core/components/foo';" > lib/core/app/components/foo.js` | export component from addon engine to main app                      |
| `yarn test:quick -f='<test name>'` -s                                                    | run tests in the browser, filtering by test name                    |

### Code Generators

Make use of the many generators for code, try `ember help generate` for more details. If you're using a component that can be widely-used, consider making it an `addon` component instead (see [this PR](https://github.com/hashicorp/vault/pull/6629) for more details)

eg. a reusable component named foo that you'd like in the core engine (read more about Ember engines [here](https://ember-engines.com/docs)).

- `ember g component foo --in lib/core`
- `echo "export { default } from 'core/components/foo';" > lib/core/app/components/foo.js`

### Running Tests

Running tests will spin up a Vault dev server on port :9200 via a
pretest script that testem (the test runner) executes. All of the
acceptance tests then run, which proxy requests back to that server.

- `yarn run test:oss`
- `yarn run test:oss -s` to keep the test server running after the initial run.
- `yarn run test -f="policies"` to filter the tests that are run. `-f` gets passed into
  [QUnit's `filter` config](https://api.qunitjs.com/config/QUnit.config#qunitconfigfilter-string--default-undefined)

### Linting

- `yarn lint:js`
- `yarn lint:hbs`
- `yarn lint:fix`

### Further Reading / Useful Links

- [ember.js](https://emberjs.com/)
- [ember-cli](https://cli.emberjs.com/release/)
- Development Browser Extensions
  - [ember inspector for chrome](https://chrome.google.com/webstore/detail/ember-inspector/bmdblncegkenkacieihfhpjfppoconhi)
  - [ember inspector for firefox](https://addons.mozilla.org/en-US/firefox/addon/ember-inspector/)

## Best Practices
