#!/usr/bin/env node
/* eslint-env node */

/**
 * Codemod to transform BasicDropdown and Tooltip trigger and content components
 * As of version 2 of ember-basic-dropdown the yielded component names are now capitalized
 * In addition, splattributes are used and class must be passed as an attribute rather than argument
 */

module.exports = () => {
  const evalTag = (tag, prefix) => {
    switch (tag) {
      case `${prefix}.trigger`:
        return `${prefix}.Trigger`;
      case `${prefix}.content`:
        return `${prefix}.Content`;
    }
  };
  // remove @ symbol from class and change @tagName to @htmlTag
  const transformAttrs = (attrs) => {
    attrs.forEach((attr) => {
      if (attr.name === '@class') {
        attr.name = 'class';
      } else if (attr.name.includes('tagName')) {
        attr.name = '@htmlTag';
      }
    });
  };

  return {
    ElementNode(node) {
      // BasicDropdown and Tooltip components yield variable is either d, dd, or T
      const prefixes = ['d', 'D', 'dd', 'T'];
      // capitalize trigger and content and transform attributes
      for (let i = 0; i < prefixes.length; i++) {
        const newTag = evalTag(node.tag, prefixes[i]);
        if (newTag) {
          node.tag = newTag;
          transformAttrs(node.attributes);
          break;
        }
      }
    },
  };
};
