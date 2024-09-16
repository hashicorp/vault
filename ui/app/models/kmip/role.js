/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import apiPath from 'vault/utils/api-path';
import lazyCapabilities from 'vault/macros/lazy-capabilities';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';

@withExpandedAttributes()
export default class KmipRoleModel extends Model {
  @attr({ readOnly: true }) backend;
  @attr({ readOnly: true }) scope;

  get editableFields() {
    return Object.keys(this.allByKey).filter((k) => !['backend', 'scope', 'role'].includes(k));
  }

  @lazyCapabilities(apiPath`${'backend'}/scope/${'scope'}/role/${'id'}`, 'backend', 'scope', 'id') updatePath;
}
