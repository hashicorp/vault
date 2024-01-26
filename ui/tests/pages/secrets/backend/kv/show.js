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
  deleteBtn: clickable('[data-test-secret-delete] button'),
  deleteBtnV1: clickable('[data-test-secret-v1-delete="true"] button'),
  deleteBtnV2: clickable('[data-test-secret-v2-delete="true"] button'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  rows: collection('data-test-row-label'),
  toggleJSON: clickable('[data-test-secret-json-toggle]'),
  toggleIsPresent: isPresent('[data-test-secret-json-toggle]'),
  edit: clickable('[data-test-secret-edit]'),
  editIsPresent: isPresent('[data-test-secret-edit]'),
  noReadIsPresent: isPresent('[data-test-write-without-read-empty-message]'),
  noReadMessage: text('data-test-empty-state-message'),

  deleteSecret() {
    return this.deleteBtn().confirmBtn();
  },
  deleteSecretV1() {
    return this.deleteBtnV1().confirmBtn();
  },
  deleteSecretV2() {
    return this.deleteBtnV2().confirmBtn();
  },
});
