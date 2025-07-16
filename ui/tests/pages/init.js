/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, collection, visitable, fillable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/init'),
  shares: fillable('[data-test-key-shares]'),
  threshold: fillable('[data-test-key-threshold]'),
  keys: collection('[data-test-key-box]'),
  init: async function (shares, threshold) {
    await this.visit();
    return this.shares(shares).threshold(threshold);
  },
});
