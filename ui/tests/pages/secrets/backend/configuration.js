/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable, text } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/secrets/:backend/configuration'),
  defaultTTL: text('[data-test-value-div="Default Lease TTL"]'),
  maxTTL: text('[data-test-value-div="Max Lease TTL"]'),
});
