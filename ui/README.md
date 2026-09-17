**Table of Contents**

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

# Vault UI

This README outlines the details of collaborating on this Ember application.

## Ember Version Upgrade Matrix

Respective versions for `ember-cli`, `ember-source` and `ember-data` for each version of Vault that contains an upgrade.

| Vault Version | Ember CLI | Ember Source | Ember Data |
| ------------- | --------- | ------------ | ---------- |
| 1.19.x        | 5.8.0     | 5.8.0        | 5.3.2      |
| 1.17.x        | 5.4.2     | 5.4.0        | 4.12.4     |
| 1.15.x        | 4.12.1    | 4.12.0       | 4.11.3     |
| 1.13.x        | 4.4.0     | 4.4.4        | 4.5.0      |
| 1.11.x        | 3.28.5    | 3.28.10      | 3.28.6     |
| 1.10.x        | 3.24.0    | 3.24.7       | 3.24.0     |
| 1.9.x         | 3.22.0    | 3.22.0       | 3.22.0     |

## Prerequisites

You will need the following things properly installed on your computer.

- [Git](https://git-scm.com/)
- [nvm](https://github.com/nvm-sh/nvm) to install and switch to the Node.js version mirrored in the repo root `.nvmrc` and `.node-version` files
- [pnpm](https://pnpm.io/)
- [Google Chrome](https://google.com/chrome/)

## Running a Vault Server

Before running Vault UI locally, a Vault server must be running. First, ensure
Vault dev is built according the instructions in `../README.md`.

- To start a single local Vault server: `pnpm vault`
- To start a local Vault cluster: `pnpm vault:cluster`

These commands may also be [aliased on your local device](https://github.com/hashicorp/vault-tools/blob/master/users/noelle/vault_aliases).

## Running the UI locally

To spin up the UI, a Vault server must be running (see previous step).
_All of the commands below assume you're in the `ui/` directory._

> These steps will start an Ember CLI server that proxies requests to port 8200,
> and enable live rebuilding of the application as you change the UI application code.
> Visit your app at [http://localhost:4200](http://localhost:4200).

1. Use the mirrored Node.js version from the repo root:

`nvm use`

2. Install dependencies:

`pnpm i`

3. Run Vault UI and proxy back to a Vault server running on the default port, 8200:

`pnpm start`

> If your Vault server is running on a different port you can use the
> long-form version of the npm script:

`pnpm exec ember server --proxy=http://localhost:PORT`

### Mirage

[Mirage](https://miragejs.com/docs/getting-started/introduction/) can be helpful for mocking backend endpoints.
Look in [mirage/handlers](mirage/handlers/) for existing mocked backends.

Run pnpm with mirage: `export MIRAGE_DEV_HANDLER=<handler> && pnpm start`

Where `handlername` is one of the options exported in [mirage/handlers/index](mirage/handlers/index.js)

To stop using the handler, kill the pnpm process (Ctrl+c) and then unset the environment variable.
`unset MIRAGE_DEV_HANDLER`

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
| `pnpm start`                                      | start the app with live reloading (vault must be running on port :8200) |
| `export MIRAGE_DEV_HANDLER=<handler>; pnpm start` | start the app with the mocked mirage backend, with handler provided     |
| `make static-dist && make dev-ui`                 | build a Vault binary with UI assets (run from root directory not `/ui`) |
| `pnpm exec ember g component foo -ir core`        | generate a component in the /addon engine                               |
| `pnpm test:filter`                                | run non-enterprise in the browser                                       |
| `pnpm test:filter -f='<test name>'`               | run tests in the browser, filtering by test name                        |
| `pnpm lint:js`                                    | lint javascript files                                                   |

### Code Generators

Make use of the many generators for code, try `pnpm exec ember help generate` for more details. If you're using a component that can be widely-used, consider making it an `addon` component instead (see [this PR](https://github.com/hashicorp/vault/pull/6629) for more details)

eg. a reusable component named foo that you'd like in the core engine (read more about Ember engines [here](https://ember-engines.com/docs)).

- `pnpm exec ember g component foo -ir core`

The above command creates a template-only component by default. If you'd like to add a backing class, add the `-gc` flag:

- `pnpm exec ember g component foo -gc -ir core`

### Running Tests

Running tests will spin up a Vault dev server on port :9200 via a pretest script that testem (the test runner) executes. All of the acceptance tests then run, which proxy requests back to that server. The normal test scripts use `ember-exam` which split into parallel runs, which is excellent for speed but makes it harder to debug. So we have a custom package script that automatically opens all the tests in a browser, and we can pass the `-f` flag to target the test(s) we're debugging.

- `pnpm run test` lint & run all the tests (CI uses this)
- `pnpm run test:oss` lint & run all the non-enterprise tests (CI uses this)
- `pnpm run test:quick` run all the tests without linting
- `pnpm run test:quick-oss` run all the non-enterprise tests without linting
- `pnpm run test:filter -f="policies"` run the filtered test in the browser with no splitting. `-f` is set to `!enterprise` by default
  [QUnit's `filter` config](https://api.qunitjs.com/config/QUnit.config#qunitconfigfilter-string--default-undefined)

#### Playwright personas

Playwright runs against real Vault servers with the UI embedded in the binary. Build
the binary with the intended UI first (`make static-dist && make entdev-ui` from the
repository root), and make it available as `vault` on `PATH`. Playwright starts,
initializes, and unseals its own servers; do not start servers on its test ports.

From `ui/`, run a focused spec with its persona:

```sh
pnpm exec playwright test e2e/tests/raft/storage.spec.ts --project=chrome:raft
```

[test-users.ts](e2e/test-users.ts) defines personas independently of
[policy definitions](e2e/policies/index.ts):

- `name` identifies the setup/browser projects, test directory, and session/key files.
- `storage` selects an implemented backend (`inmem` or `raft`).
- `policy` selects a key from `USER_POLICY_MAP`, not a second persona name.

The existing `superuser` and `raft` personas share the `superuser` policy, but run
on separate in-memory and Raft servers. Their project names remain
`chrome:superuser` and `chrome:raft`. Add a persona entry and a matching
`e2e/tests/<name>/` directory to add a scenario; add a policy only when permissions
actually differ. Setup verifies the backend and issued token's policy, then saves
an isolated session for the browser project.

Only implemented environment dimensions belong in the persona type. Consul,
license variants, edition, and version selection need runner support and tests
before being added; an unused field would not prove that environment was tested.
The Raft persona includes a single-node overview test and an isolated membership
flow in [membership.spec.ts](e2e/tests/raft/membership.spec.ts). The latter starts
two disposable Raft nodes on dynamically allocated loopback ports. It joins an
uninitialized peer through the UI, unseals it, verifies both backend membership and
the reloaded overview, then removes that peer through the UI and verifies its
persisted absence. It uses a non-root persona token for membership read-back and
removal. Nodes and temporary data are cleaned up even when the test fails, and
retries get a fresh cluster. Credential-entry traces, videos, and screenshots are
disabled for this test.
Non-secret membership responses and rendered rows are attached to the test report
before joining, after joining/reload, and after removal/reload.

```sh
pnpm exec playwright test e2e/tests/raft --project=chrome:raft
```

For migration verification, run the same membership spec against binaries built
with the pre- and post-migration UI and retain both results. Passing against only
one binary does not establish before/after equivalence.

### Linting

- `pnpm lint:js`
- `pnpm lint:hbs`
- `pnpm lint:fix`

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
