/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

const MANAGED_AUTH_BACKENDS = ['cert', 'userpass', 'ldap', 'okta', 'radius'];

export function supportedManagedAuthBackends() {
  return MANAGED_AUTH_BACKENDS;
}

export default buildHelper(supportedManagedAuthBackends);
