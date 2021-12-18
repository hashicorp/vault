#!/usr/bin/env node
/* eslint-env node */
/**
 * Codemod to convert args to attributes for Input and TextArea built in components
 * eg. <Input @id="foo" /> -> <Input id="foo" />
 */

module.exports = () => {
  // partial list of deprecated arguments
  // complete list used by linter found at:
  // https://github.com/ember-template-lint/ember-template-lint/blob/master/lib/rules/no-unknown-arguments-for-builtin-components.js
  const deprecatedArgs = [
    '@id',
    '@name',
    '@autocomplete',
    '@spellcheck',
    '@disabled',
    '@class',
    '@placeholder',
    '@wrap',
    '@rows',
    '@readonly',
    '@step',
    '@min',
    '@pattern',
  ];
  return {
    ElementNode(node) {
      if (['Textarea', 'Input', 'LinkTo'].includes(node.tag)) {
        const attrs = node.attributes;
        let i = 0;
        while (i < attrs.length) {
          const arg = deprecatedArgs.find((name) => name === attrs[i].name);
          if (arg) {
            attrs[i].name = arg.slice(1);
          }
          i++;
        }
      }
    },
  };
};
