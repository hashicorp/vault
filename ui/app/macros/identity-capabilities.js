/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default function () {
  return lazyCapabilities(apiPath`identity/${'identityType'}/id/${'id'}`, 'id', 'identityType');
}
