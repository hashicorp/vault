# Vault UI Storybook

This README outlines the details of running and collaborating on the Vault UI Storybook.

## Prerequisites

You will need to have all the dependencies for the [Vault UI](../README.md#prerequisites) installed.

## Commands at a Glance

| Command                                    | Description               |
| ------------------------------------------ | ------------------------- |
| `yarn`                                     | install dependencies      |
| `yarn storybook`                           | run storybook             |
| `ember generate story [name-of-component]` | generate a new story      |
| `yarn gen-story-md [name-of-component]`    | update a story notes file |

## Contributing Guide

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

2. Generate a new story with `ember generate story [name-of-component]`
3. Inside the newly generated `stories` file, add at least one example of the component. If the component should be interactive, enable the [Storybook Knobs addon](https://github.com/storybooks/storybook/tree/master/addons/knobs).
4. Generate the `notes` file for the component with run `yarn gen-story-md [name-of-component]` (e.g. `yarn gen-md alert-banner`). This will generate markdown documentation of the component and place it at `vault/ui/stories/[name-of-component].md`.

See the [Storybook Docs](https://storybook.js.org/docs/basics/introduction/) for more information on writing stories.

### Code Generators

It is important to keep the story and notes files up to date. To ease the process of creating and updating stories we use code generators.

## Resources

- [Storybook for Ember Live Example](https://storybooks-ember.netlify.com/?path=/story/addon-centered--button)
- [Storybook Addons](https://github.com/storybooks/storybook/tree/master/addons/)
- [Storybook Docs](https://storybook.js.org/docs/basics/introduction/)
