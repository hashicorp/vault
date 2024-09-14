/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import apiPath from 'vault/utils/api-path';
import lazyCapabilities from 'vault/macros/lazy-capabilities';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';
import { removeFromArray } from 'vault/helpers/remove-from-array';

@withExpandedAttributes()
export default class KmipRoleModel extends Model {
  @attr({ readOnly: true }) backend;
  @attr({ readOnly: true }) scope;
  @attr({ readOnly: true }) name;

  editableFields() {
    return removeFromArray(Object.keys(this.allByKey), ['backend', 'scope', 'name']);
  }

  @lazyCapabilities(apiPath`${'backend'}/scope/${'scope'}/role/${'id'}`, 'backend', 'scope', 'id') updatePath;
}
