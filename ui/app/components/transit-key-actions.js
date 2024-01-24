/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { assert } from '@ember/debug';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { encodeString } from 'vault/utils/b64';

/**
 * @module TransitKeyActions
 * TransitKeyActions component handles the actions a user can take on a transit key model. The model and props are updated on every tab change
 *
 * @example
 * ```js
 * <TransitKeyActions
 * @key={{this.model}}
 * @selectedAction="hmac"
 * />
 *
 * @param {string} selectedAction - This is the query param "action" value. Ex: hmac, verify, decrypt, etc.
 */

const STARTING_TRANSIT_PARAMS = {
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
export default class TransitKeyActions extends Component {
  @service store;
  @service flashMessages;
  @service router;

  @tracked isModalActive = false;
  @tracked errors = null;
  @tracked props = Object.assign({}, STARTING_TRANSIT_PARAMS); // shallow copy of the object. We don't want to mutate the original.

  constructor() {
    super(...arguments);
    assert(`@selectedAction is required for TransitKeyActions components`, this.args.selectedAction);
    assert('@key` is required for TransitKeyActions components', this.args.key);

    if (this.args.selectedAction === 'export') {
      this.props.exportKeyType = this.args.key.exportKeyTypes.firstObject;
      this.props.exportKeyVersion = this.args.key.validKeyVersions.lastObject;
    }
  }
  @action updateProps() {
    // reset props and errors to null. this is called when the queryParam changes, i.e. the tab is changed.
    this.errors = null; // reset errors
    this.props = Object.assign({}, STARTING_TRANSIT_PARAMS); // reset props
  }

  get keyIsRSA() {
    const { type } = this.args.key;
    return type === 'rsa-2048' || type === 'rsa-3072' || type === 'rsa-4096';
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

  @task
  @waitFor
  *doSubmit(data, options = {}, maybeEvent) {
    const event = options.type === 'submit' ? options : maybeEvent;
    if (event) {
      event.preventDefault();
    }
    const { backend, id } = this.args.key;
    const action = this.args.selectedAction;
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
      this.errors = e.errors;
    }
  }

  handleSuccess(resp, options, action) {
    if (resp && resp.data) {
      if (action === 'export' && resp.data.keys) {
        const { keys, type, name } = resp.data;
        resp.data.keys = { keys, type, name };
      }
      this.props = { ...this.props, ...resp.data };
    }
    if (options.wrapTTL) {
      this.props = { ...this.props, ...{ wrappedToken: resp.wrap_info.token } };
    }
    // open the modal
    this.isModalActive = !this.isModalActive;
    // verify doesn't trigger a success message
    if (this.selectedAction !== 'verify') {
      this.flashMessages.success(SUCCESS_MESSAGE_FOR_ACTION[action]);
    }
  }
}
