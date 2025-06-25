/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, fillable, visitable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets/:backend/list'),
  visitShow: visitable('/vault/secrets/:backend/show/:id'),
  visitCreate: visitable('/vault/secrets/:backend/create'),
  name: fillable('[data-test-input="name"]'),
  type: fillable('[data-test-input="type"'),
  tweakSource: fillable('[data-test-input="tweak_source"'),
  maskingChar: fillable('[data-test-input="masking_character"'),
});
