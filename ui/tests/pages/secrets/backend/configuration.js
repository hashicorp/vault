/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable, text } from 'ember-cli-page-object';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

export default create({
  visit: visitable('/vault/secrets/:backend/configuration'),
  defaultTTL: text(GENERAL.infoRowValue('Default Lease TTL')),
  maxTTL: text(GENERAL.infoRowValue('Max Lease TTL')),
});
