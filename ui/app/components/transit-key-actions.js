import { assign } from '@ember/polyfills';
import { copy } from 'ember-copy';
import { assert } from '@ember/debug';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { set, get, computed } from '@ember/object';

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
export default Component.extend(TRANSIT_PARAMS, {
  store: service(),

  // public attrs
  selectedAction: null,
  key: null,

  onRefresh() {},
  init() {
    this._super(...arguments);
    if (get(this, 'selectedAction')) {
      return;
    }
    set(this, 'selectedAction', get(this, 'key.supportedActions.firstObject'));
    assert('`key` is required for `' + this.toString() + '`.', this.getModelInfo());
  },

  didReceiveAttrs() {
    this._super(...arguments);
    this.checkAction();
    if (get(this, 'selectedAction') === 'export') {
      this.setExportKeyDefaults();
    }
  },

  setExportKeyDefaults() {
    const exportKeyType = get(this, 'key.exportKeyTypes.firstObject');
    const exportKeyVersion = get(this, 'key.validKeyVersions.lastObject');
    this.setProperties({
      exportKeyType,
      exportKeyVersion,
    });
  },

  keyIsRSA: computed('key.type', function() {
    let type = get(this, 'key.type');
    return type === 'rsa-2048' || type === 'rsa-4096';
  }),

  getModelInfo() {
    const model = get(this, 'key') || get(this, 'backend');
    if (!model) {
      return null;
    }
    const backend = get(model, 'backend') || get(model, 'id');
    const id = get(model, 'id');

    return {
      backend,
      id,
    };
  },

  checkAction() {
    const currentAction = get(this, 'selectedAction');
    const oldAction = get(this, 'oldSelectedAction');

    this.resetParams(oldAction, currentAction);
    set(this, 'oldSelectedAction', currentAction);
  },

  resetParams(oldAction, action) {
    let params = copy(TRANSIT_PARAMS);
    let paramsToKeep;
    let clearWithoutCheck =
      !oldAction ||
      // don't save values from datakey
      oldAction === 'datakey' ||
      // can rewrap signatures — using that as a ciphertext later would be problematic
      (oldAction === 'rewrap' && !get(this, 'key.supportsEncryption'));

    if (!clearWithoutCheck && action) {
      paramsToKeep = PARAMS_FOR_ACTION[action];
    }

    if (paramsToKeep) {
      paramsToKeep.forEach(param => delete params[param]);
    }
    //resets params still left in the object to defaults
    this.setProperties(params);
    if (action === 'export') {
      this.setExportKeyDefaults();
    }
  },

  handleError(e) {
    this.set('errors', e.errors);
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
    this.setProperties(props);
    if (action === 'rotate') {
      this.get('onRefresh')();
    }
  },

  compactData(data) {
    let type = get(this, 'key.type');
    let isRSA = type === 'rsa-2048' || type === 'rsa-4096';
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
      this.resetParams(null, get(this, 'selectedAction'));
    },

    clearParams(params) {
      const arr = Array.isArray(params) ? params : [params];
      arr.forEach(param => this.set(param, null));
    },

    doSubmit(data, options = {}) {
      const { backend, id } = this.getModelInfo();
      const action = this.get('selectedAction');
      let payload = data ? this.compactData(data) : null;
      this.setProperties({
        errors: null,
        result: null,
      });
      this.get('store')
        .adapterFor('transit-key')
        .keyAction(action, { backend, id, payload }, options)
        .then(
          resp => this.handleSuccess(resp, options, action),
          (...errArgs) => this.handleError(...errArgs)
        );
    },
  },
});
