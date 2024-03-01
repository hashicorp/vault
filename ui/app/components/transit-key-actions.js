/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { assert } from '@ember/debug';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { encodeString } from 'vault/utils/b64';
import errorMessage from 'vault/utils/error-message';

/**
 * @module TransitKeyActions
 * TransitKeyActions component handles the actions a user can take on a transit key model. The model and props are updated on every tab change
 *
 * @example
 * <TransitKeyActions
 * @key={{this.model}}
 * @selectedAction="hmac"
 * />
 *
 * @param {string} selectedAction - This is the query param "action" value. Ex: hmac, verify, decrypt, etc. The only time this param can be empty is if a user is exporting a key
 * @param {object} key - This is the transit key model.
 */

const STARTING_TRANSIT_PROPS = {
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
  wrappedTTL: '30m',
  valid: null,
  plaintextOriginal: null,
  didDecode: false,
  verification: 'Signature',
};

const PROPS_TO_KEEP = {
  encrypt: ['plaintext', 'context', 'nonce', 'key_version'],
  decrypt: ['ciphertext', 'context', 'nonce'],
  sign: ['input', 'hash_algorithm', 'key_version', 'prehashed', 'signature_algorithm'],
  verify: ['input', 'hmac', 'signature', 'hash_algorithm', 'prehashed'],
  hmac: ['input', 'algorithm', 'key_version'],
  rewrap: ['ciphertext', 'context', 'nonce', 'key_version'],
  datakey: [],
};

const SUCCESS_MESSAGE_FOR_ACTION = {
  sign: 'Signed your data.',
  // the verify action doesn't trigger a success message
  hmac: 'Created your hash output.',
  encrypt: 'Created a wrapped token for your data.',
  decrypt: 'Decrypted the data from your token.',
  rewrap: 'Created a new token for your data.',
  datakey: 'Generated your key.',
  export: 'Exported your key.',
};

export default class TransitKeyActions extends Component {
  @service store;
  @service flashMessages;
  @service router;

  @tracked isModalActive = false;
  @tracked errors = null;
  @tracked props = { ...STARTING_TRANSIT_PROPS }; // Shallow copy of the object. We don't want to mutate the original.

  constructor() {
    super(...arguments);
    assert('@key is required for TransitKeyActions components', this.args.key);

    if (this.firstSupportedAction === 'export' || this.args.selectedAction === 'export') {
      this.props.exportKeyType = this.args.key.exportKeyTypes[0];
      this.props.exportKeyVersion = this.args.key.validKeyVersions.lastObject;
    }
  }

  get keyIsRSA() {
    const { type } = this.args.key;
    return type === 'rsa-2048' || type === 'rsa-3072' || type === 'rsa-4096';
  }

  get firstSupportedAction() {
    return this.args.key.supportedActions[0];
  }

  handleSuccess(resp, options, action) {
    if (resp && resp.data) {
      if (action === 'export' && resp.data.keys) {
        const { keys, type, name } = resp.data;
        resp.data.keys = { keys, type, name };
      }
      this.props = { ...this.props, ...resp.data };

      // While we do not pass ciphertext as a value to the JsonEditor, so that navigating from rewrap to decrypt will not show ciphertext in the editor, we still want to clear it from the props after rewrapping.
      if (action === 'rewrap' && !this.args.key.supportsEncryption) {
        this.props.ciphertext = null;
      }
    }
    if (options.wrapTTL) {
      this.props = { ...this.props, wrappedToken: resp.wrap_info.token };
    }
    this.isModalActive = true;
    // verify doesn't trigger a success message
    if (this.args.selectedAction !== 'verify') {
      this.flashMessages.success(SUCCESS_MESSAGE_FOR_ACTION[action]);
    }
  }

  @action updateProps() {
    this.errors = null;
    // There are specific props we want to carry over from the previous tab.
    // Ex: carrying over this.props.context from the encrypt tab to the decrypt tab, but not carrying over this.props.plaintext.
    // To do this, we make a new object that contains the old this.props key/values from the previous tab that we want to keep. We then merge that new object into the STARTING_TRANSIT_PROPS object to come up with our new this.props tracked property.
    // This action is passed to did-update in the component.
    const transferredProps = PROPS_TO_KEEP[this.args.selectedAction]?.reduce(
      (obj, key) => ({ ...obj, [key]: this.props[key] }),
      {}
    );
    this.props = { ...STARTING_TRANSIT_PROPS, ...transferredProps };
  }

  compactData(data) {
    return Object.keys(data).reduce((result, key) => {
      if (key === 'signature_algorithm' && !this.keyIsRSA) {
        return result;
      }
      if (data[key]) {
        result[key] = data[key];
      }
      return result;
    }, {});
  }

  @action toggleEncodeBase64() {
    this.props.encodedBase64 = !this.props.encodedBase64;
  }

  @action clearSpecificProps(arrayToClear) {
    arrayToClear.forEach((prop) => (this.props[prop] = null));
  }

  @task
  @waitFor
  *doSubmit(data, options = {}, maybeEvent) {
    this.errors = null;
    const event = options.type === 'submit' ? options : maybeEvent;
    if (event) {
      event.preventDefault();
    }
    const { backend, id } = this.args.key;
    const action = this.args.selectedAction || this.firstSupportedAction;
    const { ...formData } = data || {};
    if (!this.props.encodedBase64) {
      if (action === 'encrypt' && !!formData.plaintext) {
        formData.plaintext = encodeString(formData.plaintext);
      }
      if ((action === 'hmac' || action === 'verify' || action === 'sign') && !!formData.input) {
        formData.input = encodeString(formData.input);
      }
    }
    const payload = formData ? this.compactData(formData) : null;

    try {
      const resp = yield this.store
        .adapterFor('transit-key')
        .keyAction(action, { backend, id, payload }, options);

      this.handleSuccess(resp, options, action);
    } catch (e) {
      this.errors = errorMessage(e);
    }
  }
}
