/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class OidcScopeSerializer extends ApplicationSerializer {
  primaryKey = 'name';
}
