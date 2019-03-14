# Vault UI Storybook

This README outlines the details of running and collaborating on the Vault UI Storybook.

## Prerequisites

You will need to have all the dependencies for the [Vault UI](../README.md#prerequisites) installed.

## Commands at a Glance

| Command | Description |
|---|---|
| `yarn` | install dependencies |
| `yarn storybook` | run storybook |
| `ember generate story [name-of-component]` | generate a new story |
| TBA | update a story notes file |

## Contributing Guide

### Writing Stories

For the most part, each component in `vault/ui/app/components` should have a corresponding `[component-name].stories.js` and `[component-name].md` files within `vault/ui/stories`. 

#### Which components should have stories?

- Components that are heavily reused

#### Which shouldn't?

- Components which are overly complex
- Components that are only used once

#### Adding a new story

1. Generate a new story with `ember generate story [name-of-component]`
2. Inside the newly generated `stories` file, add at least one example of the component. If the component should be interactive, enable the [Storybook Knobs addon](https://github.com/storybooks/storybook/tree/master/addons/knobs).
3. Inside the newly generated `[component].md` file, fill out the table of component properties, with a Usage example and a link to the component source code.

See the [Storybook Docs](https://storybook.js.org/docs/basics/introduction/) for more information on writing stories.

### Code Generators

It is  important to keep the story and notes files up to date. To ease the process of creating and updating stories we use code generators.

## Resources

- [Storybook for Ember Live Example](https://storybooks-ember.netlify.com/?path=/story/addon-centered--button)
- [Storybook Addons](https://github.com/storybooks/storybook/tree/master/addons/)
- [Storybook Docs](https://storybook.js.org/docs/basics/introduction/)
