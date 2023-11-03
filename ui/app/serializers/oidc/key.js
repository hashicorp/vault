/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class OidcKeySerializer extends ApplicationSerializer {
  primaryKey = 'name';
}
