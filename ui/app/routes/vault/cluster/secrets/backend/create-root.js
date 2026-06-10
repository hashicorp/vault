/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import { service } from '@ember/service';
import EditBase from './secret-edit';
import KeymgmtKeyForm from 'vault/forms/keymgmt/key';
import KeymgmtProviderForm from 'vault/forms/keymgmt/provider';
import TotpKeyForm from 'vault/forms/totp/key';
import TransitKeyForm from 'vault/forms/transit/key';
import SshRoleForm from 'vault/forms/ssh/role';
import AlphabetForm from 'vault/forms/transform/alphabet';
import TemplateForm from 'vault/forms/transform/template';
import RoleForm from 'vault/forms/transform/role';
import TransformationForm from 'vault/forms/transform/transformation';
import { KeyManagementUpdateKeyRequestTypeEnum } from '@hashicorp/vault-client-typescript';

const secretModel = (store, backend, key) => {
  const model = store.createRecord('secret', {
    path: key,
  });
  return model;
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

    if (modelType === 'totp-key') {
      return new TotpKeyForm(
        {
          backend,
          generate: true,
          algorithm: 'SHA1',
          digits: 6,
          period: 30,
          exported: true,
          key_size: 20,
          skew: 1,
          qr_size: 200,
        },
        { isNew: true }
      );
    }

    if (modelType === 'transit-key') {
      return new TransitKeyForm(
        {
          backend,
          type: 'aes256-gcm96',
          auto_rotate_period: '0s',
        },
        { isNew: true }
      );
    }

    if (modelType === 'role-ssh') {
      return new SshRoleForm(
        { backend, key_type: 'ca', not_before_duration: '30s', port: 22 },
        { isNew: true }
      );
    }
    if (modelType === 'transform/alphabet') {
      return new AlphabetForm({ backend }, { isNew: true });
    }
    if (modelType === 'transform/template') {
      return new TemplateForm({ backend }, { isNew: true });
    }
    if (modelType === 'transform/role') {
      return new RoleForm({ backend }, { isNew: true });
    }
    if (modelType === 'transform') {
      return new TransformationForm({ backend, type: 'fpe' }, { isNew: true });
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
