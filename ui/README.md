# vault

This README outlines the details of collaborating on this Ember application.
A short introduction of this app could easily go here.

## Prerequisites

You will need the following things properly installed on your computer.

- [Node.js](https://nodejs.org/) (with NPM)
- [Yarn](https://yarnpkg.com/en/)
- [Git](https://git-scm.com/)
- [Ember CLI](https://ember-cli.com/)
- [Husky\*](https://github.com/typicode/husky)
- [lint-staged\*](https://www.npmjs.com/package/lint-staged)

\* Husky and lint-staged are optional dependencies - running `yarn` will install them.
If don't want them installed (husky adds files for every hooks in `.git/hooks/`),
then you can run `yarn --ignore-optional`. If you've ignored the optional deps
previously and want to install them, you have to tell yarn to refetch all deps by
running `yarn --force`.

## Running / Development

To get all of the JavaScript dependencies installed, run this in the `ui` directory:

`yarn`

If you want to run Vault UI and proxy back to a Vault server running
on the default port, 8200, run the following in the `ui` directory:

- `yarn run start`

This will start an Ember CLI server that proxies requests to port 8200,
and enable live rebuilding of the application as you change the UI application code.
Visit your app at [http://localhost:4200](http://localhost:4200).

If your Vault server is running on a different port you can use the
long-form version of the npm script:

`ember server --proxy=http://localhost:PORT`

### Code Generators

Make use of the many generators for code, try `ember help generate` for more details

### Running Tests

Running tests will spin up a Vault dev server on port 9200 via a
pretest script that testem (the test runner) executes. All of the
acceptance tests then run, proxing requests back to that server.

- `yarn run test-oss`
- `yarn run test-oss -s` to keep the test server running after the initial run.
- `yarn run test -f="policies"` to filter the tests that are run. `-f` gets passed into
  [QUnit's `filter` config](https://api.qunitjs.com/config/QUnit.config#qunitconfigfilter-string--default-undefined)

### Linting

- `yarn lint:hbs`
- `yarn lint:js`
- `yarn lint:js -- --fix`

### Building Vault UI into a Vault Binary

We use `go-bindata-assetfs` to build the static assets of the
Ember application into a Vault binary.

This can be done by running these commands from the root directory run:
`make static-dist`
`make dev-ui`

This will result in a Vault binary that has the UI built-in - though in
a non-dev setup it will still need to be enabled via the `ui` config or
setting `VAULT_UI` environment variable.

## Further Reading / Useful Links

- [ember.js](http://emberjs.com/)
- [ember-cli](https://ember-cli.com/)
- Development Browser Extensions
  - [ember inspector for chrome](https://chrome.google.com/webstore/detail/ember-inspector/bmdblncegkenkacieihfhpjfppoconhi)
  - [ember inspector for firefox](https://addons.mozilla.org/en-US/firefox/addon/ember-inspector/)
