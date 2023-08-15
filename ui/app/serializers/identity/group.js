/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { EmbeddedRecordsMixin } from '@ember-data/serializer/rest';
import IdentitySerializer from './_base';

export default IdentitySerializer.extend(EmbeddedRecordsMixin, {
  attrs: {
    alias: { embedded: 'always' },
  },

  normalizeFindRecordResponse(store, primaryModelClass, payload) {
    if (payload.alias && Object.keys(payload.alias).length === 0) {
      delete payload.alias;
    }
    return this._super(...arguments);
  },

  serialize() {
    const json = this._super(...arguments);
    delete json.alias;
    if (json.type === 'external') {
      delete json.member_entity_ids;
      delete json.member_group_ids;
    }
    return json;
  },
});
