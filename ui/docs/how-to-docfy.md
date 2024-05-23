---
title: How to docfy
order: 1
---

# How to docfy

http://localhost:4200/ui/docs (navigate directly or click the present icon in the bottom right of the app footer)

## `docs/` markdown files

Side nav links correspond to the file + directory structure within the `docs/` directory. These markdown files can be edited directly and any updates should be saved and pushed to main.

## generating component docs

The `docs/components/` directory is where generated markdown files for components will live after running `yarn docs`. These files are included `.gitignore` so they are not pushed to main. If `jsdoc-to-markdown` errors it will be printed in the console.

| Command                             | Description                                                                |
| ----------------------------------- | -------------------------------------------------------------------------- |
| `yarn docs`                         | generate markdown file for every\* component in the `addon/core` directory |
| `yarn docfy-md some-component-name` | generate markdown file for specific component                              |
| `yarn docfy-md read-more core`      | generate markdown for `read-more` component in the `core` addon            |
| `rm -f ./docs/components/*`         | cleanup and delete generated component markdown files                      |

> \*replication components are skipped as these are technically not reused outside of the replication engine and should not live in the addon engine

## Writing documentation for a component

1. Accurate `jsdoc` syntax is important so `jsdoc-to-markdown` properly generates the markdown file for that component.

2. Docfy renders an actual instance of the component beneath `@example` as a sample. Make sure component uses proper hbs syntax. The component args **cannot span multiple lines**.

### jsdoc examples:

- **confirmation-modal** [github link](https://github.com/hashicorp/vault/blob/main/ui/lib/core/addon/components/confirmation-modal.js) | [VScode link](../lib/core/addon/components/confirmation-modal.js)
- **certificate-card** [github link](https://github.com/hashicorp/vault/blob/main/ui/lib/core/addon/components/certificate-card.js) | [VScode link](../lib/core/addon/components/certificate-card.js)

### @deprecated example

- **alert-inline** [github link](https://github.com/hashicorp/vault/blob/main/ui/lib/core/addon/components/alert-inline.js) | [VScode link](../lib/core/addon/components/alert-inline.js)

### Syntax tips

> - Param types: `object`, `string`, `function`, `array`
> - Do not include `null` for empty default values
> - The script automatically wraps default string values in quotes, do not include them in the jsdoc default values
> - Wrap optional params in brackets `[]`

```
/**
 * @module ComponentName
 * Description of the component
 *
 * @example
 <ComponentName @param={{}} optionalParam={{}} />
 *
 * @param {type} paramName - description
 * @param {string} requiredParam=foo - Do not wrap default values in strings
 * @param {array} [optionalParamName] - An optional parameter
 * @param {string} [param=some default value] - An optional parameter with a default value
 */
```

3. Check the markdown file for syntax errors or typos, then navigate to the component url `http://localhost:4200/ui/docs/components/some-component-name`

4. Fix the jsdoc and rerun `yarn docfy-md some-component-name` to check.

### More info

- [Building and consuming components](./building-components.md)
