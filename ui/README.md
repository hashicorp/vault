** Table of Contents **

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Vault UI](#vault-ui)
  - [Ember CLI Version Matrix](#ember-cli-version-matrix)
  - [Prerequisites](#prerequisites)
  - [Running a Vault Server](#running-a-vault-server)
  - [Running / Development](#running--development)
    - [Code Generators](#code-generators)
    - [Running Tests](#running-tests)
    - [Linting](#linting)
    - [Building Vault UI into a Vault Binary](#building-vault-ui-into-a-vault-binary)
  - [Vault Storybook](#vault-storybook)
    - [Storybook Commands at a Glance](#storybook-commands-at-a-glance)
    - [Writing Stories](#writing-stories)
      - [Adding a new story](#adding-a-new-story)
    - [Code Generators](#code-generators-1)
    - [Storybook Deployment](#storybook-deployment)
  - [Further Reading / Useful Links](#further-reading--useful-links)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Vault UI

This README outlines the details of collaborating on this Ember application.

## Ember CLI Version Matrix

| Vault Version | Ember Version |
| ------------- | ------------- |
| 1.10.x        | 3.28.5        |
| 1.9.x         | 3.22.0        |
| 1.8.x         | 3.22.0        |
| 1.7.x         | 3.11          |

## Prerequisites

You will need the following things properly installed on your computer.

- [Node.js](https://nodejs.org/) (with NPM)
- [Yarn](https://yarnpkg.com/en/)
- [Git](https://git-scm.com/)
- [Ember CLI](https://ember-cli.com/)
- [lint-staged\*](https://www.npmjs.com/package/lint-staged)

\* lint-staged is an optional dependency - running `yarn` will install it.
If don't want optional dependencies installed you can run `yarn --ignore-optional`. If you've ignored the optional deps
previously and want to install them, you have to tell yarn to refetch all deps by
running `yarn --force`.

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

## Running / Development

To get all of the JavaScript dependencies installed, run this in the `ui` directory:

- `yarn`

If you want to run Vault UI and proxy back to a Vault server running
on the default port, 8200, run the following in the `ui` directory:

- `yarn start`

This will start an Ember CLI server that proxies requests to port 8200,
and enable live rebuilding of the application as you change the UI application code.
Visit your app at [http://localhost:4200](http://localhost:4200).

If your Vault server is running on a different port you can use the
long-form version of the npm script:

`ember server --proxy=http://localhost:PORT`

To run yarn with mirage, do:

- `yarn start:mirage handlername`

Where `handlername` is one of the options exported in `mirage/handlers/index`

### Code Generators

Make use of the many generators for code, try `ember help generate` for more details. If you're using a component that can be widely-used, consider making it an `addon` component instead (see [this PR](https://github.com/hashicorp/vault/pull/6629) for more details)

eg. a reusable component named foo that you'd like in the core engine

- `ember g component foo --in lib/core`
- `echo "export { default } from 'core/components/foo';" > lib/core/app/components/foo.js`

### Running Tests

Running tests will spin up a Vault dev server on port 9200 via a
pretest script that testem (the test runner) executes. All of the
acceptance tests then run, proxing requests back to that server.

- `yarn run test:oss`
- `yarn run test:oss -s` to keep the test server running after the initial run.
- `yarn run test -f="policies"` to filter the tests that are run. `-f` gets passed into
  [QUnit's `filter` config](https://api.qunitjs.com/config/QUnit.config#qunitconfigfilter-string--default-undefined)

### Linting

- `yarn lint`
- `yarn lint:fix`

### Building Vault UI into a Vault Binary

We use the [embed](https://golang.org/pkg/embed/) package from Go 1.16+ to build
the static assets of the Ember application into a Vault binary.

This can be done by running these commands from the root directory run:
`make static-dist`
`make dev-ui`

This will result in a Vault binary that has the UI built-in - though in
a non-dev setup it will still need to be enabled via the `ui` config or
setting `VAULT_UI` environment variable.

## Vault Storybook

The Vault UI uses Storybook to catalog all of its components. Below are details for running and contributing to Storybook.

### Storybook Commands at a Glance

| Command                                                                  | Description                                                |
| ------------------------------------------------------------------------ | ---------------------------------------------------------- |
| `yarn storybook`                                                         | run storybook                                              |
| `ember generate story [name-of-component]`                               | generate a new story                                       |
| `ember generate story [name-of-component] -ir [name-of-engine-or-addon]` | generate a new story in the specified engine or addon      |
| `yarn gen-story-md [name-of-component]`                                  | update a story notes file                                  |
| `yarn gen-story-md [name-of-component] [name-of-engine-or-addon]`        | update a story notes file in the specified engine or addon |

### Writing Stories

Each component in `vault/ui/app/components` should have a corresponding `[component-name].stories.js` and `[component-name].md` files within `vault/ui/stories`. Components in the `core` addon located at `vault/ui/lib/core/addon/components` have corresponding stories and markdown files in `vault/ui/lib/core/stories`.

#### Adding a new story

1. Make sure the component is well-documented using [jsdoc](http://usejsdoc.org/tags-exports.html). This documentation should at minimum include the module name, an example of usage, and the params passed into the handlebars template. For example, here is how we document the ToggleButton Component:

````js
/**
 * @module ToggleButton
 * `ToggleButton` components are used to expand and collapse content with a toggle.
 *
 * @example
 * ```js
 *   <ToggleButton @openLabel="Encrypt Output with PGP" @closedLabel="Encrypt Output with PGP" @toggleTarget={{this}} @toggleAttr="showOptions"/>
 *  {{#if showOptions}}
 *     <div>
 *       <p>
 *         I will be toggled!
 *       </p>
 *     </div>
 *   {{/if}}
 * ```
 *
 * @param {String} toggleAttr=null - The attribute upon which to toggle.
 * @param {Object} attrTarget=null - The target upon which the event handler should be added.
 * @param {String} [openLabel=Hide options] - The message to display when the toggle is open. //optional params are denoted by square brackets
 * @param {String} [closedLabel=More options] - The message to display when the toggle is closed.
 */
````

Note that placing a param inside brackets (e.g. `[closedLabel=More options]` indicates it is optional and has a default value of `'More options'`.)

2. Generate a new story with `ember generate story [name-of-component]`
3. Inside the newly generated `stories` file, add at least one example of the component. If the component should be interactive, enable the [Storybook Knobs addon](https://github.com/storybooks/storybook/tree/master/addons/knobs).
4. Generate the `notes` file for the component with `yarn gen-story-md [name-of-component] [name-of-engine-or-addon]` (e.g. `yarn gen-md alert-banner core`). This will generate markdown documentation of the component and place it at `vault/ui/stories/[name-of-component].md`. If your component is a template-only component, you will need to manually create the markdown file. The markdown file will need to be imported in your `[component-name].stories.js` file (e.g. `import notes from './[name-of-component].md'`).
5. The completed `[component-name].stories.js` file should look something like this (with knobs):

```js
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { text, withKnobs } from '@storybook/addon-knobs';
import notes from './stat-text.md';

storiesOf('MyComponent', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `MyComponent`,
    () => ({
      template: hbs`
      <h5 class="title is-5">My Component</h5>
      <MyComponent @param={{param}} @anotherParam={{anotherParam}} />
    `,
      context: {
        param: text('param', 'My parameter'),
        anotherParam: boolean('anotherParam', true),
      },
    }),
    { notes }
  );
```

See the [Storybook Docs](https://storybook.js.org/docs/basics/introduction/) for more information on writing stories.

### Code Generators

It is important to add all new components into Storybook and to keep the story and notes files up to date. To ease the process of creating and updating stories please use the code generators using the [commands listed above](#storybook-commands-at-a-glance).

### Storybook Deployment

A Vercel integration deploys a static Storybook build for any PR on the Vault GitHub repo. A preview link will show up in the PR checks. Once items are merged, the auto-deployed integration will publish that build making it available at [https://vault-storybook.vercel.app](https://vault-storybook.vercel.app). Currently the Vercel integration will cd into the `ui/` directory and then run `yarn deploy:storybook` so troubleshooting any issues can be done locally by running this same command. The logs for this build are public and will be linked from the PR checks.

## Further Reading / Useful Links

- [ember.js](http://emberjs.com/)
- [ember-cli](https://ember-cli.com/)
- Development Browser Extensions
  - [ember inspector for chrome](https://chrome.google.com/webstore/detail/ember-inspector/bmdblncegkenkacieihfhpjfppoconhi)
  - [ember inspector for firefox](https://addons.mozilla.org/en-US/firefox/addon/ember-inspector/)
- [Storybook for Ember Live Example](https://vault-storybook.vercel.app/?path=/story/addon-centered--button)
- [Storybook Addons](https://github.com/storybooks/storybook/tree/master/addons/)
- [Storybook Docs](https://storybook.js.org/docs/basics/introduction/)
