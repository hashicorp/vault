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
  // We have to leave this here because the backend doesn't support the file type yet.
  credentials: attr('string', {
    editType: 'file',
  }),

  googleCertsEndpoint: attr('string'),

  fieldGroups: computed('newFields', function () {
    let groups = [
      { default: ['credentials'] },
      {
        'Google Cloud Options': ['googleCertsEndpoint'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }
    return fieldToAttrs(this, groups);
  }),
});
