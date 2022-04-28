import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action, setProperties } from '@ember/object';
import { copy } from 'ember-copy';
import { assign } from '@ember/polyfills';
import { encodeString } from 'vault/utils/b64';

/**
 * @module TransitKeyActionsTwo
 * TransitKeyActionsTwo components are used to...
 *
 * @example
 * ```js
 * <TransitKeyActionsTwo @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {string} key - requiredParam is...
 * @param {string} selectedAction - optionalParam is...
 * @param {string} [key=null] - param1 is...
 * @param {object} [capabilities] - param1 is...
 * @param {function} [onRefresh] - param1 is...
 * @param {string} [backend] - param1 is...
 */

export const TRANSIT_PARAMS = {
  hash_algorithm: 'sha2-256',
  algorithm: 'sha2-256',
  signature_algorithm: 'pss',
  bits: 256,
  bytes: 32,
  ciphertext: null,
  context: null, // ARG TODO probably remove
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
export const PARAMS_FOR_ACTION = {
  sign: ['input', 'hash_algorithm', 'key_version', 'prehashed', 'signature_algorithm'],
  verify: ['input', 'hmac', 'signature', 'hash_algorithm', 'prehashed'],
  hmac: ['input', 'algorithm', 'key_version'],
  encrypt: ['plaintext', 'context', 'nonce', 'key_version'],
  decrypt: ['ciphertext', 'context', 'nonce'],
  rewrap: ['ciphertext', 'context', 'nonce', 'key_version'],
};
export const SUCCESS_MESSAGE_FOR_ACTION = {
  sign: 'Signed your data',
  // the verify action doesn't trigger a success message
  hmac: 'Created your hash output',
  encrypt: 'Created a wrapped token for your data',
  decrypt: 'Decrypted the data from your token',
  rewrap: 'Created a new token for your data',
  datakey: 'Generated your key',
  export: 'Exported your key',
};
export default class TransitKeyActionsTwo extends Component {
  @service store;
  @service flashMessages;

  @tracked context = null;
  @tracked exportKeyType;
  @tracked exportKeyVersion;
  @tracked errors;
  @tracked isModalActive = false;
  @tracked oldSelectedAction;
  @tracked ciphertext;
  @tracked plaintext;
  @tracked key_version;

  constructor() {
    super(...arguments);

    if (this.args.selectedAction) {
      return;
    }
    console.log('run test suite and see if this ever happens');
    // debugger;
  }

  setExportKeyDefaults() {
    const exportKeyType = this.args.key.exportKeyTypes.firstObject;
    const exportKeyVersion = this.args.key.validKeyVersions.lastObject;
    setProperties(this, {
      exportKeyType,
      exportKeyVersion,
    });
  }

  get keyIsRSA() {
    let type = this.args.key.type;
    return type === 'rsa-2048' || type === 'rsa-3072' || type === 'rsa-4096';
  }

  getModelInfo() {
    const model = this.args.key || this.args.backend;
    if (!model) {
      return null;
    }
    const backend = model.backend || model.id;
    const id = model.id;

    return {
      backend,
      id,
    };
  }

  checkAction() {
    const currentAction = this.args.selectedAction;
    this.resetParams(this.oldSelectedAction, currentAction); // ARG TODO ???
    this.oldSelectedAction = this.args.selectedAction;
  }

  resetParams(oldAction, action) {
    let params = copy(TRANSIT_PARAMS);
    let paramsToKeep;
    let clearWithoutCheck =
      !oldAction ||
      // don't save values from datakey
      oldAction === 'datakey' ||
      // can rewrap signatures — using that as a ciphertext later would be problematic
      (oldAction === 'rewrap' && !this.args.key.supportsEncryption);

    if (!clearWithoutCheck && action) {
      paramsToKeep = PARAMS_FOR_ACTION[action];
    }

    if (paramsToKeep) {
      paramsToKeep.forEach((param) => delete params[param]);
    }
    // resets params still left in the object to defaults
    this.clearErrors();
    // debugger;
    setProperties(this, params);
    if (action === 'export') {
      this.setExportKeyDefaults();
    }
  }

  handleError(e) {
    this.errors = e.errors;
  }

  clearErrors() {
    this.errors = null;
  }

  triggerSuccessMessage(action) {
    const message = SUCCESS_MESSAGE_FOR_ACTION[action];
    if (!message) return;
    this.flashMessages.success(message);
  }

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
      this.isModalActive = !this.isModalActive; // toggle
      setProperties(this, props);
    }
    if (action === 'rotate') {
      this.onRefresh();
    }
    this.triggerSuccessMessage(action);
  }

  compactData(data) {
    let type = this.args.key.type;
    let isRSA = type === 'rsa-2048' || type === 'rsa-3072' || type === 'rsa-4096';
    return Object.keys(data).reduce((result, key) => {
      if (key === 'signature_algorithm' && !isRSA) {
        return result;
      }
      if (data[key]) {
        result[key] = data[key];
      }
      return result;
    }, {});
  }

  @action
  onActionChange(action) {
    this.args.selectedAction = action; // ARG TODO not going to work
    this.checkAction();
  }

  @action
  onClear() {
    this.resetParams(null, this.args.selectedAction); // ARG TODO not going to work.
  }

  @action
  clearParams(params) {
    const arr = Array.isArray(params) ? params : [params];
    arr.forEach((param) => (param ? null : null)); // linting won't let me do param = null
  }

  @action
  toggleModal(successMessage) {
    if (!!successMessage && typeof successMessage === 'string') {
      this.flashMessages.success(successMessage);
    }
    this.isModalActive = !this.isModalActive; // toggle
  }

  @action
  doSubmit(data, options = {}) {
    const { backend, id } = this.getModelInfo();
    const action = this.args.selectedAction;
    const { encodedBase64, ...formData } = data || {};
    if (!encodedBase64) {
      if (action === 'encrypt' && !!formData.plaintext) {
        formData.plaintext = encodeString(formData.plaintext);
      }
      if ((action === 'hmac' || action === 'verify' || action === 'sign') && !!formData.input) {
        formData.input = encodeString(formData.input);
      }
    }
    let payload = formData ? this.compactData(formData) : null;
    setProperties(this, {
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
  }

  @action
  onHash(newValue) {
    this.context = newValue;
  }
}
