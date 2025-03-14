/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Base } from '../show';
import { create, clickable, collection } from 'ember-cli-page-object';

export default create({
  ...Base,
  rows: collection('data-test-row-label'),
  deleteBtn: clickable('[data-test-confirm-action-trigger]'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  deleteKey() {
    return this.deleteBtn().confirmBtn();
  },
});
