#!/usr/bin/env node
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-env node */

// print to console all files that include LinkTo elements using the {{on modifier}}
module.exports = (env) => {
  let fileAlerted;
  return {
    ElementNode(node) {
      if (node.tag === 'LinkTo') {
        if (!fileAlerted) {
          const usesModifier = node.modifiers.find((modifier) => modifier.path.original === 'on');
          if (usesModifier) {
            console.log(env.filePath); // eslint-disable-line
            fileAlerted = true;
          }
        }
      }
    },
  };
};
