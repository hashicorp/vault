#!/usr/bin/env node
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-env node */

/**
 * Codemod to convert quoteless attribute or argument to mustache statement
 * eg. data-test-foo=true -> data-test-foo={{true}}
 * eg @isVisible=true -> @isVisible={{true}}
 */

module.exports = (env) => {
  const { builders } = env.syntax;
  return {
    ElementNode({ attributes }) {
      let i = 0;
      while (i < attributes.length) {
        const { type, chars } = attributes[i].value;
        if (type === 'TextNode' && chars && !attributes[i].quoteType) {
          attributes[i].value = builders.mustache(builders.path(attributes[i].value.chars));
        }
        i++;
      }
    },
  };
};
