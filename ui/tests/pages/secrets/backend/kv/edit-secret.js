/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Base } from '../create';
import { clickable, create, fillable } from 'ember-cli-page-object';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

export default create({
  ...Base,
  path: fillable('[data-test-secret-path="true"]'),
  secretKey: fillable('[data-test-secret-key]'),
  secretValue: fillable('[data-test-secret-value] textarea'),
  save: clickable(GENERAL.saveButton),
  toggleJSON: clickable('[data-test-toggle-input="json"]'),
});
