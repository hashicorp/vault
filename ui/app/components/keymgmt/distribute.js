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
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class KeymgmtDistribute extends Component {
  @service store;

  @tracked keyModel;
  @tracked providerType;

  distributeKey(backend, kms, key) {
    let adapter = this.store.adapterFor('keymgmt/key');
    return adapter.distribute(backend, kms, key);
  }

  async createKey(keyName) {
    const key = await this.store.createRecord(`keymgmt/key`, {
      backend: this.args.backend,
      id: keyName,
    });
    this.keyModel = key;
  }

  destroyKey() {
    if (this.keyModel) {
      console.log(this.keyModel);
      this.keyModel
        .destroyRecord()
        .then(() => {
          this.keyModel = null;
        })
        .catch((e) => {
          console.log('could not destroy record');
          console.log({ e });
        });
      // record.destroyRecord();
      // this.keyModel = null;
    }
  }
  async setProviderType(id) {
    if (!id) {
      this.providerType = '';
      return;
    }

    if (id === 'example-kms') {
      this.providerType = 'gcpckms';
    } else {
      this.providerType = 'azurekeyvault';
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

  get operations() {
    const pt = this.providerType;
    if (pt === 'awskms') {
      return ['encrypt', 'decrypt'];
    } else if (pt === 'gcpckms') {
      const kt = this.keyModel?.type || ''; // TODO: How do we store & retrieve from existing key
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
    this.args.model.set('provider', evt.target.value);
    if (evt.target.value) {
      this.setProviderType(evt.target.value);
    }
  }
  @action
  handleKeyType(evt) {
    this.keyModel.set('type', evt.target.value);
  }

  @action
  handleOperation(evt) {
    const ops = [...this.args.model.operations];
    if (evt.target.checked) {
      ops.push(evt.target.id);
    } else {
      const idx = ops.indexOf(evt.target.id);
      ops.splice(idx, 1);
    }
    this.args.model.set('operations', ops);
  }

  @action
  async handleKeySelect(selected) {
    const selectedKey = selected[0] || selected;
    if (!selectedKey) {
      this.args.model.set('key', null);
      this.destroyKey();
      return;
    } else if (selectedKey.isNew) {
      this.createKey(selectedKey.id);
    }
    this.args.model.set('key', selectedKey.id);
  }

  @action
  createDistribution(evt) {
    console.log({ evt });

    evt.preventDefault();
    // clean checkboxes against provider type
    // const { backend } = this.args.model;
    // TODO: Check if we need to create a key first
    // this.distributeKey(backend, 'example-kms', 'example-key')
    //   .then(() => {
    //     console.log('success');
    //   })
    //   .catch((e) => {
    //     console.log('error', e);
    //   });
  }
}
