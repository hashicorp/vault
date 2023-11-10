---
title: How to doc
order: 1
---

# How to write component docs

1. Write properly formatted `jsdoc` for the component. The component under `@example` needs the accurate syntax so docfy can render an actual example of the component. If the component args span multiple lines, do not add an asterisk at the beginning of each line.

> _Syntax notes:_
>
> - Param types: `object`, `string`, `function`, `array`
> - Do not include include `null` for empty default values
> - The script automatically wraps default string values in quotes, do not include them in the jsdoc default values

```
/**
 * @module ComponentName
 * Description of the component
 *
 * @example
 <ComponentName @param={{}} optionalParam={{}} />
 *
 * @param {type} paramName - description
 * @param {array} [optionalParamName] - An optional parameter
 * @param {string} [param=some default value] - An option parameter with a default value
 */
```

2. Run `$ yarn docfy-md some-component-name` to generate the markdown file. It will add it to the `docs/components/` directory. If the the component is in an add-on or separate ember engine, include the name of engine. For example, if a component lives in the `core` addon run:
   `yarn docfy-md some-component-name core`

3. Check the markdown file for syntax errors or typos, then navigate to the component url `http://localhost:4200/ui/docs/components/some-component-name`

### More info

- [Building and consuming components](./building-components.md)
