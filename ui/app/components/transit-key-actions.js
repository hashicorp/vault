import { assign } from '@ember/polyfills';
import { copy } from 'ember-copy';
import { assert } from '@ember/debug';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { set, get, computed } from '@ember/object';
import { encodeString } from 'vault/utils/b64';

const TRANSIT_PARAMS = {
  hash_algorithm: 'sha2-256',
  algorithm: 'sha2-256',
  signature_algorithm: 'pss',
  bits: 256,
  bytes: 32,
  ciphertext: null,
  context: null,
  format: 'base64',
  hmac: null,
  input: null,
  key_version: 0,
  keys: null,
  nonce: null,
  param: 'wrapped',
  prehashed: false,
  plaintext: null,
  random_bytes: null,
  signature: null,
  sum: null,
  encodedBase64: false,
  exportKeyType: null,
  exportKeyVersion: null,
  wrappedToken: null,
  valid: null,
  plaintextOriginal: null,
  didDecode: false,
  verification: 'Signature',
};
const PARAMS_FOR_ACTION = {
  sign: ['input', 'hash_algorithm', 'key_version', 'prehashed', 'signature_algorithm'],
  verify: ['input', 'hmac', 'signature', 'hash_algorithm', 'prehashed'],
  hmac: ['input', 'algorithm', 'key_version'],
  encrypt: ['plaintext', 'context', 'nonce', 'key_version'],
  decrypt: ['ciphertext', 'context', 'nonce'],
  rewrap: ['ciphertext', 'context', 'nonce', 'key_version'],
};
const SUCCESS_MESSAGE_FOR_ACTION = {
  sign: 'Signed your data',
  // the verify action doesn't trigger a success message
  hmac: 'Created your hash output',
  encrypt: 'Created a wrapped token for your data',
  decrypt: 'Decrypted the data from your token',
  rewrap: 'Created a new token for your data',
  datakey: 'Generated your key',
  export: 'Exported your key',
};
export default Component.extend(TRANSIT_PARAMS, {
  store: service(),
  flashMessages: service(),

  // public attrs
  selectedAction: null,
  key: null,
  isModalActive: false,

  onRefresh() {},
  init() {
    this._super(...arguments);
    // eslint-disable-next-line ember/no-get
    if (this.selectedAction) {
      return;
    }
    // eslint-disable-next-line ember/no-get
    set(this, 'selectedAction', get(this, 'key.supportedActions.firstObject'));
    assert('`key` is required for `' + this.toString() + '`.', this.getModelInfo());
  },

  didReceiveAttrs() {
    this._super(...arguments);
    this.checkAction();
    if (this.selectedAction === 'export') {
      this.setExportKeyDefaults();
    }
  },

  setExportKeyDefaults() {
    const exportKeyType = this.key.exportKeyTypes.firstObject;
    const exportKeyVersion = this.key.validKeyVersions.lastObject;
    this.setProperties({
      exportKeyType,
      exportKeyVersion,
    });
  },

  keyIsRSA: computed('key.type', function () {
    const type = this.key.type;
    return type === 'rsa-2048' || type === 'rsa-3072' || type === 'rsa-4096';
  }),

  getModelInfo() {
    const model = this.key || this.backend;
    if (!model) {
      return null;
    }
    const backend = model.backend || model.id;
    const id = model.id;

    return {
      backend,
      id,
    };
  },

  checkAction() {
    const currentAction = this.selectedAction;
    const oldAction = this.oldSelectedAction;

    this.resetParams(oldAction, currentAction);
    set(this, 'oldSelectedAction', currentAction);
  },

  resetParams(oldAction, action) {
    const params = copy(TRANSIT_PARAMS);
    let paramsToKeep;
    const clearWithoutCheck =
      !oldAction ||
      // don't save values from datakey
      oldAction === 'datakey' ||
      // can rewrap signatures — using that as a ciphertext later would be problematic
      (oldAction === 'rewrap' && !this.key.supportsEncryption);

    if (!clearWithoutCheck && action) {
      paramsToKeep = PARAMS_FOR_ACTION[action];
    }

    if (paramsToKeep) {
      paramsToKeep.forEach((param) => delete params[param]);
    }
    //resets params still left in the object to defaults
    this.clearErrors();
    this.setProperties(params);
    if (action === 'export') {
      this.setExportKeyDefaults();
    }
  },

  handleError(e) {
    this.set('errors', e.errors);
  },

  clearErrors() {
    this.set('errors', null);
  },

  triggerSuccessMessage(action) {
    const message = SUCCESS_MESSAGE_FOR_ACTION[action];
    if (!message) return;
    this.flashMessages.success(message);
  },

  handleSuccess(resp, options, action) {
    let props = {};
    if (resp && resp.data) {
      if (action === 'export' && resp.data.keys) {
        const { keys, type, name } = resp.data;
        resp.data.keys = { keys, type, name };
      }
      props = assign({}, props, resp.data);
    }
    if (options.wrapTTL) {
      props = assign({}, props, { wrappedToken: resp.wrap_info.token });
    }
    if (!this.isDestroyed && !this.isDestroying) {
      this.toggleProperty('isModalActive');
      this.setProperties(props);
    }
    if (action === 'rotate') {
      this.onRefresh();
    }
    this.triggerSuccessMessage(action);
  },

  compactData(data) {
    const type = this.key.type;
    const isRSA = type === 'rsa-2048' || type === 'rsa-3072' || type === 'rsa-4096';
    return Object.keys(data).reduce((result, key) => {
      if (key === 'signature_algorithm' && !isRSA) {
        return result;
      }
      if (data[key]) {
        result[key] = data[key];
      }
      return result;
    }, {});
  },

  actions: {
    onActionChange(action) {
      set(this, 'selectedAction', action);
      this.checkAction();
    },

    onClear() {
      this.resetParams(null, this.selectedAction);
    },

    clearParams(params) {
      const arr = Array.isArray(params) ? params : [params];
      arr.forEach((param) => this.set(param, null));
    },

    toggleModal(successMessage) {
      if (!!successMessage && typeof successMessage === 'string') {
        this.flashMessages.success(successMessage);
      }
      this.toggleProperty('isModalActive');
    },

    doSubmit(data, options = {}, maybeEvent) {
      const event = options.type === 'submit' ? options : maybeEvent;
      if (event) {
        event.preventDefault();
      }
      const { backend, id } = this.getModelInfo();
      const action = this.selectedAction;
      const { encodedBase64, ...formData } = data || {};
      if (!encodedBase64) {
        if (action === 'encrypt' && !!formData.plaintext) {
          formData.plaintext = encodeString(formData.plaintext);
        }
        if ((action === 'hmac' || action === 'verify' || action === 'sign') && !!formData.input) {
          formData.input = encodeString(formData.input);
        }
      }
      const payload = formData ? this.compactData(formData) : null;
      this.setProperties({
        errors: null,
        result: null,
      });
      this.store
        .adapterFor('transit-key')
        .keyAction(action, { backend, id, payload }, options)
        .then(
          (resp) => this.handleSuccess(resp, options, action),
          (...errArgs) => this.handleError(...errArgs)
        );
    },
  },
});
