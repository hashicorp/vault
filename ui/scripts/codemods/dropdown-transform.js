#!/usr/bin/env node
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */

/**
 * Codemod to transform BasicDropdown and Tooltip trigger and content components
 * As of version 2 of ember-basic-dropdown the yielded component names are now capitalized
 * In addition, splattributes are used and class must be passed as an attribute rather than argument
 */

module.exports = () => {
  return {
    ElementNode(node) {
      // ensure we have the right parent node by first looking for BasicDropdown or ToolTip nodes
      if (['BasicDropdown', 'ToolTip'].includes(node.tag)) {
        node.children.forEach((child) => {
          // capitalize trigger and content and transform attributes
          if (child.type === 'ElementNode' && child.tag.match(/\.(content|trigger)/gi)) {
            const { tag } = child;
            const char = tag.charAt(tag.indexOf('.') + 1);
            child.tag = tag.replace(char, char.toUpperCase());
            // remove @ symbol from class and change @tagName to @htmlTag
            // Content component does not use splattributes -- apply class with @defaultClass arg
            child.attributes.forEach((attr) => {
              if (attr.name.includes('class')) {
                if (child.tag.includes('Content')) {
                  attr.name = '@defaultClass';
                } else if (attr.name === '@class') {
                  attr.name = 'class';
                }
              } else if (attr.name.includes('tagName')) {
                attr.name = '@htmlTag';
              }
            });
          }
        });
      }
    },
  };
};
