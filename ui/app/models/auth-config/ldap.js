/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { attr } from '@ember-data/model';
import { computed } from '@ember/object';

import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

export default AuthConfig.extend({
  useOpenAPI: true,
  certificate: attr({
    label: 'Certificate',
    editType: 'file',
  }),
  fieldGroups: computed('newFields', function () {
    let groups = [
      {
        default: ['url'],
      },
      {
        'LDAP Options': [
          'starttls',
          'insecureTls',
          'discoverdn',
          'denyNullBind',
          'tlsMinVersion',
          'tlsMaxVersion',
          'certificate',
          'clientTlsCert',
          'clientTlsKey',
          'userattr',
          'upndomain',
          'anonymousGroupSearch',
        ],
      },
      {
        'Customize User Search': ['binddn', 'userdn', 'bindpass', 'userfilter'],
      },
      {
        'Customize Group Membership Search': ['groupfilter', 'groupattr', 'groupdn', 'useTokenGroups'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }
    return fieldToAttrs(this, groups);
  }),
});
