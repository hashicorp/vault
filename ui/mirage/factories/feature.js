/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  feature_flags() {
    return []; // VAULT_CLOUD_ADMIN_NAMESPACE
  },
});
