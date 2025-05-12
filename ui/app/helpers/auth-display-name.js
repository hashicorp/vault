/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ALL_LOGIN_METHODS } from 'vault/utils/supported-login-methods';

export default function authDisplayName(type) {
  const displayName = ALL_LOGIN_METHODS?.find((t) => t.type === type)?.displayName;
  return displayName || type;
}
