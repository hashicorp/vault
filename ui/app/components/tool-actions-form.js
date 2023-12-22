/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { match } from '@ember/object/computed';
import { assign } from '@ember/polyfills';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { setProperties, computed, set } from '@ember/object';
import { addSeconds, parseISO } from 'date-fns';
import { A } from '@ember/array';
import { capitalize } from '@ember/string';

const DEFAULTS = {
  token: null,
  rewrap_token: null,
  errors: A(),
  wrap_info: null,
  creation_time: null,
  creation_ttl: null,
  data: '{\n}',
  unwrap_data: null,
  details: null,
  wrapTTL: null,
  sum: null,
  random_bytes: null,
  input: null,
};

const WRAPPING_ENDPOINTS = ['lookup', 'wrap', 'unwrap', 'rewrap'];

export default Component.extend(DEFAULTS, {
  flashMessages: service(),
  store: service(),
  // putting these attrs here so they don't get reset when you click back
  //random
  bytes: 32,
  //hash
  format: 'base64',
  algorithm: 'sha2-256',

  tagName: '',

  didReceiveAttrs() {
    this._super(...arguments);
    this.checkAction();
  },

  selectedAction: null,

  reset() {
    if (this.isDestroyed || this.isDestroying) {
      return;
    }
    setProperties(this, DEFAULTS);
  },

  checkAction() {
    const currentAction = this.selectedAction;
    const oldAction = this.oldSelectedAction;

    if (currentAction !== oldAction) {
      this.reset();
    }
    set(this, 'oldSelectedAction', currentAction);
  },

  dataIsEmpty: match('data', new RegExp(DEFAULTS.data)),

  expirationDate: computed('creation_time', 'creation_ttl', function () {
    const { creation_time, creation_ttl } = this;
    if (!(creation_time && creation_ttl)) {
      return null;
    }
    // returns new Date with seconds added.
    return addSeconds(parseISO(creation_time), creation_ttl);
  }),

  handleError(e) {
    set(this, 'errors', e.errors);
  },

  handleSuccess(resp, action) {
    let props = {};
    const secret = (resp && resp.data) || resp.auth;
    if (secret && action === 'unwrap') {
      const details = {
        'Request ID': resp.request_id,
        'Lease ID': resp.lease_id || 'None',
        Renewable: resp.renewable ? 'Yes' : 'No',
        'Lease Duration': resp.lease_duration || 'None',
      };
      props = assign({}, props, { unwrap_data: secret }, { details: details });
    }
    props = assign({}, props, secret);
    if (resp && resp.wrap_info) {
      const keyName = action === 'rewrap' ? 'rewrap_token' : 'token';
      props = assign({}, props, { [keyName]: resp.wrap_info.token });
    }
    setProperties(this, props);
    this.flashMessages.success(`${capitalize(action)} was successful.`);
  },

  getData() {
    const action = this.selectedAction;
    if (WRAPPING_ENDPOINTS.includes(action)) {
      return this.dataIsEmpty ? { token: (this.token || '').trim() } : JSON.parse(this.data);
    }
    if (action === 'random') {
      return { bytes: this.bytes, format: this.format };
    }
    if (action === 'hash') {
      return { input: this.input, format: this.format, algorithm: this.algorithm };
    }
  },

  actions: {
    doSubmit(evt) {
      evt.preventDefault();
      const action = this.selectedAction;
      const wrapTTL = action === 'wrap' ? this.wrapTTL : null;
      const data = this.getData();
      setProperties(this, {
        errors: null,
        wrap_info: null,
        creation_time: null,
        creation_ttl: null,
      });

      this.store
        .adapterFor('tools')
        .toolAction(action, data, { wrapTTL })
        .then(
          (resp) => this.handleSuccess(resp, action),
          (...errArgs) => this.handleError(...errArgs)
        );
    },

    onClear() {
      this.reset();
    },

    updateTtl(ttl) {
      set(this, 'wrapTTL', ttl);
    },

    codemirrorUpdated(val, hasErrors) {
      setProperties(this, {
        buttonDisabled: hasErrors,
        data: val,
      });
    },
  },
});
