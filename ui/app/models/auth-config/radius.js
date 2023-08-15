/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import AuthConfig from '../auth-config';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';

export default AuthConfig.extend({
  useOpenAPI: true,
  host: attr('string'),
  secret: attr('string'),

  fieldGroups: computed('newFields', function () {
    let groups = [
      {
        default: ['host', 'secret'],
      },
      {
        'RADIUS Options': ['port', 'nasPort', 'nasIdentifier', 'dialTimeout', 'unregisteredUserPolicies'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }

    return fieldToAttrs(this, groups);
  }),
});
