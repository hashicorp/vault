/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default function () {
  return lazyCapabilities(apiPath`identity/${'identityType'}/id/${'id'}`, 'id', 'identityType');
}
