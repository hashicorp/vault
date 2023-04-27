/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  feature_flags() {
    return []; // VAULT_CLOUD_ADMIN_NAMESPACE
  },
});
