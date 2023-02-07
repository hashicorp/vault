/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { create, visitable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/settings/secrets/configure/:backend/'),
});
