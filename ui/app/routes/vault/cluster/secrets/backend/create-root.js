/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import { inject as service } from '@ember/service';
import EditBase from './secret-edit';

const secretModel = (store, backend, key) => {
  const backendModel = store.peekRecord('secret-engine', backend);
  const modelType = backendModel.modelTypeForKV;
  if (modelType !== 'secret-v2') {
    const model = store.createRecord(modelType, {
      path: key,
    });
    return model;
  }
  const secret = store.createRecord(modelType);
  secret.set('engine', backendModel);
  const version = store.createRecord('secret-v2-version', {
    path: key,
  });
  secret.set('selectedVersion', version);
  return secret;
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
    if (modelType !== 'secret' && modelType !== 'secret-v2') {
      return this.store.createRecord(modelType);
    }
    // create record in capabilities that checks for access to create metadata
    // this record is then maybeQueryRecord in the component secret-create-or-update
    if (modelType === 'secret-v2') {
      // only check for kv2 secrets
      this.store.findRecord('capabilities', `${backend}/metadata/`);
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
