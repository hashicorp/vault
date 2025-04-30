/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  name: (i) => `Login rule ${i}`,
  namespace: (i) => `namespace-${i}`,
  default_auth_type: 'okta',
  backup_auth_types: () => ['oidc', 'token'],
  disable_inheritance: false,
});
