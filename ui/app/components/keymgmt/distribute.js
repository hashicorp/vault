import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { KEY_TYPES } from '../../models/keymgmt/key';

/**
 * @module KeymgmtDistribute
 * KeymgmtDistribute components are used to...
 *
 * @example
 * ```js
 * <KeymgmtDistribute @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {string} backend - name of backend, which will be the basis of other store queries
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

class DistributionData {
  @tracked key;
  @tracked provider;
  @tracked operations;
  @tracked protection;
}

const VALID_TYPES_BY_PROVIDER = {
  gcpckms: ['aes256-gcm96', 'rsa-2048', 'rsa-3072', 'rsa-4096', 'ecdsa-p256', 'ecdsa-p384', 'ecdsa-p521'],
  awskms: ['aes256-gcm96'],
  azurekeyvault: ['rsa-2048', 'rsa-3072', 'rsa-4096'],
};
export default class KeymgmtDistribute extends Component {
  @service store;

  @tracked keyModel;
  @tracked isNewKey = false;
  @tracked providerType;
  @tracked formData;

  constructor() {
    super(...arguments);
    this.formData = new DistributionData();
    this.formData.key = this.args.key || '';
    this.formData.provider = this.args.provider || '';
    if (this.args.provider) {
      this.getProviderType(this.args.provider);
    }
    if (this.args.key) {
      this.getKeyInfo(this.args.key);
    }
    this.formData.operations = [];
  }

  distributeKey(backend, kms, key) {
    let adapter = this.store.adapterFor('keymgmt/key');
    return adapter.distribute(backend, kms, key);
    // TODO: After effects
  }

  async getKeyInfo(keyName, isNew = false) {
    let key;
    if (isNew) {
      this.isNewKey = true;
      key = await this.store.createRecord(`keymgmt/key`, {
        backend: this.args.backend,
        id: keyName,
      });
    } else {
      key = await this.store.queryRecord(`keymgmt/key`, {
        backend: this.args.backend,
        id: keyName,
        recordOnly: true,
      });
    }
    console.log(key.type, 'KEY TYPE');
    this.keyModel = key;
  }

  destroyKey() {
    if (this.isNewKey) {
      // Delete key record from store if it was created here
      this.keyModel
        .destroyRecord()
        .then(() => {
          this.keyModel = null;
        })
        .catch((e) => {
          console.log('could not destroy record');
          console.log({ e });
        });
    }
    this.isNewKey = false;
    this.keyModel = null;
  }

  async getProviderType(id) {
    if (!id) {
      this.providerType = '';
      return;
    }

    if (id === 'kms-gcp') {
      this.providerType = 'gcpckms';
    } else if (id === 'kms-azure') {
      this.providerType = 'azurekeyvault';
    } else {
      this.providerType = 'awskms';
    }
    // TODO: Add back once provider model available
    // const provider = await this.store.queryRecord('keymgmt/provider', {
    //   backend: this.args.backend,
    //   id
    // });
    // this.providerType = provider.type
  }

  get keyTypes() {
    // TODO: filter these if provider type available?
    return KEY_TYPES;
  }

  get validMatchError() {
    if (!this.providerType || !this.keyModel?.type) {
      return null;
    }
    return VALID_TYPES_BY_PROVIDER[this.providerType].includes(this.keyModel.type)
      ? null
      : `${this.keyModel.type} is not a valid key for provider of type ${this.providerType}`;
  }

  get operations() {
    const pt = this.providerType;
    if (pt === 'awskms') {
      return ['encrypt', 'decrypt'];
    } else if (pt === 'gcpckms') {
      const kt = this.keyModel?.type || '';
      switch (kt) {
        case 'aes256-gcm96':
          return ['encrypt', 'decrypt'];
        case 'rsa-2048':
        case 'rsa-3072':
        case 'rsa-4096':
          return ['decrypt', 'sign'];
        case 'ecdsa-p256':
        case 'ecdsa-p384':
          return ['sign'];
        default:
          return ['encrypt', 'decrypt', 'sign', 'verify', 'wrap', 'unwrap'];
      }
    }

    return ['encrypt', 'decrypt', 'sign', 'verify', 'wrap', 'unwrap'];
  }

  @action
  handleProvider(evt) {
    this.formData.provider = evt.target.value;
    if (evt.target.value) {
      this.getProviderType(evt.target.value);
    }
  }
  @action
  handleKeyType(evt) {
    this.keyModel.set('type', evt.target.value);
  }

  @action
  handleOperation(evt) {
    const ops = [...this.formData.operations];
    if (evt.target.checked) {
      ops.push(evt.target.id);
    } else {
      const idx = ops.indexOf(evt.target.id);
      ops.splice(idx, 1);
    }
    this.formData.operations = ops;
  }

  @action
  async handleKeySelect(selected) {
    const selectedKey = selected[0] || null;
    if (!selectedKey) {
      this.formData.key = null;
      return this.destroyKey();
    }
    this.formData.key = selectedKey.id;
    return this.getKeyInfo(selectedKey.id, selectedKey.isNew);
  }

  @action
  createDistribution(evt) {
    evt.preventDefault();
    const { backend } = this.args;
    // TODO: serialize formData
    console.log(this.formData.protection, backend);
    if (this.isNewKey) {
      // TODO: Create new key before distributing
    }
    // this.distributeKey(backend, 'example-kms', 'example-key')
    //   .then(() => {
    //     console.log('success');
    //   })
    //   .catch((e) => {
    //     console.log('error', e);
    //   });
  }
}
