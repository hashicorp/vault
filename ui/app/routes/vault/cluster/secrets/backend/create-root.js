/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import { service } from '@ember/service';
import EditBase from './secret-edit';
import KeymgmtKeyForm from 'vault/forms/keymgmt/key';
import KeymgmtProviderForm from 'vault/forms/keymgmt/provider';
import { KeyManagementUpdateKeyRequestTypeEnum } from '@hashicorp/vault-client-typescript';

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

    // Handle keymgmt/key with Form class
    if (modelType === 'keymgmt/key') {
      const defaultValues = {
        backend,
        type: KeyManagementUpdateKeyRequestTypeEnum.RSA_2048,
        deletion_allowed: false,
      };
      return new KeymgmtKeyForm(defaultValues, { isNew: true });
    }

    if (modelType === 'keymgmt/provider') {
      const defaultValues = {
        backend,
        credentials: {},
      };
      return new KeymgmtProviderForm(defaultValues, { isNew: true });
    }

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
      return this.store.createRecord(modelType, { backend });
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
