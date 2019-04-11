<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Vault UI](#vault-ui)
  - [Prerequisites](#prerequisites)
  - [Running / Development](#running--development)
    - [Code Generators](#code-generators)
    - [Running Tests](#running-tests)
    - [Linting](#linting)
    - [Building Vault UI into a Vault Binary](#building-vault-ui-into-a-vault-binary)
  - [Vault Storybook](#vault-storybook)
    - [Storybook Commands at a Glance](#storybook-commands-at-a-glance)
    - [Writing Stories](#writing-stories)
      - [Adding a new story](#adding-a-new-story)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Vault UI

This README outlines the details of collaborating on this Ember application.

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

- `yarn`

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

## Vault Storybook

The Vault UI uses Storybook to catalog all of its components. Below are details for running and contributing to Storybook.

### Storybook Commands at a Glance

| Command                                    | Description               |
| ------------------------------------------ | ------------------------- |
| `yarn storybook`                           | run storybook             |
| `ember generate story [name-of-component]` | generate a new story      |
| `yarn gen-story-md [name-of-component]`    | update a story notes file |

### Writing Stories

Each component in `vault/ui/app/components` should have a corresponding `[component-name].stories.js` and `[component-name].md` files within `vault/ui/stories`.

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
 * @param toggleAttr=null {String} - The attribute upon which to toggle.
 * @param attrTarget=null {Object} - The target upon which the event handler should be added.
 * @param [openLabel=Hide options] {String} - The message to display when the toggle is open. //optional params are denoted by square brackets
 * @param [closedLabel=More options] {String} - The message to display when the toggle is closed.
 */
````
Note that placing a param inside brackets (e.g. `[closedLabel=More options]` indicates it is optional and has a default value of `'More options'`.)

2. Generate a new story with `ember generate story [name-of-component]`
3. Inside the newly generated `stories` file, add at least one example of the component. If the component should be interactive, enable the [Storybook Knobs addon](https://github.com/storybooks/storybook/tree/master/addons/knobs).
4. Generate the `notes` file for the component with `yarn gen-story-md [name-of-component]` (e.g. `yarn gen-md alert-banner`). This will generate markdown documentation of the component and place it at `vault/ui/stories/[name-of-component].md`. If your component is a template-only component, you will need to manually create the markdown file.

See the [Storybook Docs](https://storybook.js.org/docs/basics/introduction/) for more information on writing stories.

### Code Generators

It is important to add all new components into Storybook and to keep the story and notes files up to date. To ease the process of creating and updating stories please use the code generators using the [commands listed above](#storybook-commands-at-a-glance).


## Further Reading / Useful Links

- [ember.js](http://emberjs.com/)
- [ember-cli](https://ember-cli.com/)
- Development Browser Extensions
  - [ember inspector for chrome](https://chrome.google.com/webstore/detail/ember-inspector/bmdblncegkenkacieihfhpjfppoconhi)
  - [ember inspector for firefox](https://addons.mozilla.org/en-US/firefox/addon/ember-inspector/)
- [Storybook for Ember Live Example](https://storybooks-ember.netlify.com/?path=/story/addon-centered--button)
- [Storybook Addons](https://github.com/storybooks/storybook/tree/master/addons/)
- [Storybook Docs](https://storybook.js.org/docs/basics/introduction/)
