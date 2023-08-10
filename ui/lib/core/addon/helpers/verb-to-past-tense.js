/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

// only works for regular verbs. Irregular verbs ex: arise, awake, be.
export function verbToPastTense(verb, tense) {
  if (!tense) {
    assert('You must provide a tense.');
  }
  if (tense === 'gerund') {
    // ending in 'ing ex: delete => deleting || destroy => destroying
    return (
      verb
        .replace(/([^aeiouy])y$/, '$1i')
        .replace(/([^aeiouy][aeiou])([^aeiouy])$/, '$1$2$2')
        .replace(/e$/, '') + 'ing'
    );
  }
  if (tense === 'past') {
    // ending in 'ing ex: delete => deleted || destroy => destroyed
    return (
      verb
        .replace(/([^aeiouy])y$/, '$1i')
        .replace(/([^aeiouy][aeiou])([^aeiouy])$/, '$1$2$2')
        .replace(/e$/, '') + 'ed'
    );
  }
}

export default buildHelper(verbToPastTense);
