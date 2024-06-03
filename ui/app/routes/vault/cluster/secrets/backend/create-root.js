/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import { service } from '@ember/service';
import EditBase from './secret-edit';

const secretModel = (store, backend, key) => {
  const model = store.createRecord('secret', {
    path: key,
  });
  return model;
};

const transformModel = (queryParams) => {
  const modelType = 'transform';
  if (!queryParams || !queryParams.itemType) return modelType;

  return `${modelType}/${queryParams.itemType}`;
};

export default EditBase.extend({
  store: service(),

  createModel(transition) {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    let modelType = this.modelType(backend, null, { queryParams: transition.to.queryParams });
    if (modelType === 'role-ssh') {
      return this.store.createRecord(modelType, { keyType: 'ca' });
    }
    if (modelType === 'transform') {
      modelType = transformModel(transition.to.queryParams);
    }
    if (modelType === 'database/connection' && transition.to?.queryParams?.itemType === 'role') {
      modelType = 'database/role';
    }
    if (modelType !== 'secret') {
      return this.store.createRecord(modelType);
    }
    return secretModel(this.store, backend, transition.to.queryParams.initialKey);
  },

  model(params, transition) {
    return hash({
      secret: this.createModel(transition),
      capabilities: {},
    });
  },
});
