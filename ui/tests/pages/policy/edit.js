/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { clickable, create, isPresent, visitable } from 'ember-cli-page-object';
export default create({
  visit: visitable('/vault/policy/:type/:name/edit'),
  deleteIsPresent: isPresent('[data-test-confirm-action-trigger]'),
  toggleEdit: clickable('[data-test-policy-edit-toggle]'),
});
