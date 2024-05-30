/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { belongsTo } from '@ember-data/model';
import { computed } from '@ember/object';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';

export default Model.extend({
  useOpenAPI: true,
  ca: belongsTo('kmip/ca', { async: false, inverse: 'config' }),
  getHelpUrl(path) {
    return `/v1/${path}/config?help=1`;
  },

  fieldGroups: computed('newFields', function () {
    let groups = [{ default: ['listenAddrs', 'connectionTimeout'] }];

    groups = combineFieldGroups(groups, this.newFields, []);
    return fieldToAttrs(this, groups);
  }),
});
