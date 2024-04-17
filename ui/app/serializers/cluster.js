/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RESTSerializer, { EmbeddedRecordsMixin } from '@ember-data/serializer/rest';
import { decamelize } from '@ember/string';
import IdentityManager from '../utils/identity-manager';

const uuids = new IdentityManager();

export default RESTSerializer.extend(EmbeddedRecordsMixin, {
  keyForAttribute: function (attr) {
    return decamelize(attr);
  },

  attrs: {
    nodes: { embedded: 'always' },
    dr: { embedded: 'always' },
    performance: { embedded: 'always' },
  },

  setReplicationId(data) {
    if (data) {
      data.id = data.cluster_id || uuids.fetch();
    }
  },

  normalize(modelClass, data) {
    // embedded records need a unique value to be stored
    // set id for dr and performance to cluster_id or random unique id
    this.setReplicationId(data.dr);
    this.setReplicationId(data.performance);
    return this._super(modelClass, data);
  },

  pushPayload(store, payload) {
    const transformedPayload = this.normalizeResponse(
      store,
      store.modelFor('cluster'),
      payload,
      null,
      'findAll'
    );
    return store.push(transformedPayload);
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    // FIXME when multiple clusters lands
    const transformedPayload = {
      clusters: Object.assign({ id: '1' }, payload.data || payload),
    };

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
