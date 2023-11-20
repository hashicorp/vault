/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Base } from '../show';
import { create, clickable, collection, isPresent, text } from 'ember-cli-page-object';

export default create({
  ...Base,
  breadcrumbs: collection('[data-test-secret-breadcrumb]', {
    text: text(),
  }),
  deleteBtnV1: clickable('[data-test-secret-v1-delete]'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  rows: collection('data-test-row-label'),
  edit: clickable('[data-test-secret-edit]'),
  editIsPresent: isPresent('[data-test-secret-edit]'),

  deleteSecretV1() {
    return this.deleteBtnV1().confirmBtn();
  },
});
