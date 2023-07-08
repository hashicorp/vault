/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import { gt } from '@ember/object/computed';
import { camelize } from '@ember/string';
import Component from '@ember/component';
import { get, computed } from '@ember/object';
import layout from '../templates/components/shamir-flow';
import { A } from '@ember/array';

const pgpKeyFileDefault = () => ({ value: '' });
const DEFAULTS = {
  key: null,
  loading: false,
  errors: A(),
  threshold: null,
  progress: null,
  pgp_key: null,
  haveSavedPGPKey: false,
  started: false,
  generateWithPGP: false,
  pgpKeyFile: pgpKeyFileDefault(),
  nonce: '',
};

export default Component.extend(DEFAULTS, {
  tagName: '',
  store: service(),
  formText: null,
  fetchOnInit: false,
  buttonText: 'Submit',
  thresholdPath: 'required',
  generateAction: false,
  layout,

  init() {
    this._super(...arguments);
    if (this.fetchOnInit) {
      this.attemptProgress();
    }
  },

  didInsertElement() {
    this._super(...arguments);
    this.onUpdate(this.getProperties(Object.keys(DEFAULTS)));
  },

  onUpdate() {},
  onLicenseError() {},
  onShamirSuccess() {},
  // can be overridden w/an attr
  isComplete(data) {
    return data.complete === true;
  },

  stopLoading() {
    this.setProperties({
      loading: false,
      errors: [],
      key: null,
    });
  },

  reset() {
    this.setProperties(DEFAULTS);
  },

  hasProgress: gt('progress', 0),

  actionSuccess(resp) {
    const { onUpdate, isComplete, onShamirSuccess, thresholdPath } = this;
    const threshold = get(resp, thresholdPath);
    const props = {
      ...resp,
      threshold,
    };
    this.stopLoading();
    // if we have an OTP, but update doesn't include one,
    // we don't want to null it out
    if (this.otp && !props.otp) {
      delete props.otp;
    }
    this.setProperties(props);
    onUpdate(props);
    if (isComplete(props)) {
      this.reset();
      onShamirSuccess(props);
    }
  },

  actionError(e) {
    this.stopLoading();
    if (e.httpStatus === 400) {
      this.set('errors', e.errors);
    } else {
      // if licensing error, trigger parent method to handle
      if (e.httpStatus === 500 && e.errors?.join(' ').includes('licensing is in an invalid state')) {
        this.onLicenseError();
      }
      throw e;
    }
  },

  generateStep: computed('generateWithPGP', 'haveSavedPGPKey', 'pgp_key', function () {
    const { generateWithPGP, pgp_key, haveSavedPGPKey } = this;
    if (!generateWithPGP && !pgp_key) {
      return 'chooseMethod';
    }
    if (generateWithPGP) {
      if (pgp_key && haveSavedPGPKey) {
        return 'beginGenerationWithPGP';
      } else {
        return 'providePGPKey';
      }
    }
    return '';
  }),

  extractData(data) {
    const isGenerate = this.generateAction;
    const hasStarted = this.started;
    const usePGP = this.generateWithPGP;
    const nonce = this.nonce;

    if (!isGenerate || hasStarted) {
      if (nonce) {
        data.nonce = nonce;
      }
      return data;
    }

    if (usePGP) {
      return {
        pgp_key: data.pgp_key,
      };
    }

    return {
      attempt: data.attempt,
    };
  },

  attemptProgress(data) {
    const checkStatus = data ? false : true;
    let action = this.action;
    action = action && camelize(action);
    this.set('loading', true);
    const adapter = this.store.adapterFor('cluster');
    const method = adapter[action];

    method.call(adapter, data, { checkStatus }).then(
      (resp) => this.actionSuccess(resp),
      (...args) => this.actionError(...args)
    );
  },

  actions: {
    reset() {
      this.reset();
      this.set('encoded_token', null);
      this.set('otp', null);
    },

    onSubmit(data) {
      if (!data.key) {
        return;
      }
      this.attemptProgress(this.extractData(data));
    },

    startGenerate(data) {
      if (this.generateAction) {
        data.attempt = true;
      }
      this.attemptProgress(this.extractData(data));
    },

    setKey(_, keyFile) {
      this.set('pgp_key', keyFile.value);
      this.set('pgpKeyFile', keyFile);
    },

    savePGPKey() {
      if (this.pgp_key) {
        this.set('haveSavedPGPKey', true);
      }
    },
  },
});
