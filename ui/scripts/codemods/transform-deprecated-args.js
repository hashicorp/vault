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
    '@disabled', // not deprecated for LinkTo
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
      if (['Textarea', 'Input', 'LinkTo', 'ToolbarSecretLink', 'SecretLink'].includes(node.tag)) {
        const attrs = node.attributes;
        let i = 0;
        while (i < attrs.length) {
          // LinkTo uses disabled as named arg
          // ensure that it is not present as attribute since link will still work
          // since ToolbarSecretLink wraps SecretLink and SecretLink wraps LinkTo include them as well
          const disabledAsArg = ['LinkTo', 'SecretLink', 'ToolbarSecretLink'].includes(node.tag);
          if (disabledAsArg && attrs[i].name === 'disabled') {
            attrs[i].name = '@disabled';
          }
          const arg = deprecatedArgs.find((name) => {
            return node.tag.includes('SecretLink') || (node.tag === 'LinkTo' && name === '@disabled')
              ? false
              : name === attrs[i].name;
          });
          if (arg) {
            attrs[i].name = arg.slice(1);
          }
          i++;
        }
      }
    },
  };
};
