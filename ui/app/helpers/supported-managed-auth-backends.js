/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

// The UI supports management of these auth methods (i.e. configuring roles or users)
// otherwise only configuration of the method is supported.
const MANAGED_AUTH_BACKENDS = ['cert', 'kubernetes', 'ldap', 'okta', 'radius', 'userpass'];

export function supportedManagedAuthBackends() {
  return MANAGED_AUTH_BACKENDS;
}

export default buildHelper(supportedManagedAuthBackends);
