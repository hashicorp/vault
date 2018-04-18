import Ember from 'ember';

const { Component, inject, computed, get } = Ember;
const { camelize } = Ember.String;

const DEFAULTS = {
  key: null,
  loading: false,
  errors: [],
  threshold: null,
  progress: null,
  started: false,
  generateWithPGP: false,
  pgpKeyFile: { value: '' },
  nonce: '',
};

export default Component.extend(DEFAULTS, {
  tagName: '',
  store: inject.service(),
  formText: null,
  fetchOnInit: false,
  buttonText: 'Submit',
  thresholdPath: 'required',
  generateAction: false,
  encoded_token: null,

  init() {
    if (this.get('fetchOnInit')) {
      this.attemptProgress();
    }
    return this._super(...arguments);
  },

  onShamirSuccess: _ => _,
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

  hasProgress: computed.gt('progress', 0),

  actionSuccess(resp) {
    const { isComplete, onShamirSuccess, thresholdPath } = this.getProperties(
      'isComplete',
      'onShamirSuccess',
      'thresholdPath'
    );
    this.stopLoading();
    this.set('threshold', get(resp, thresholdPath));
    this.setProperties(resp);
    if (isComplete(resp)) {
      this.reset();
      onShamirSuccess(resp);
    }
  },

  actionError(e) {
    this.stopLoading();
    if (e.httpStatus === 400) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },

  extractData(data) {
    const isGenerate = this.get('generateAction');
    const hasStarted = this.get('started');
    const usePGP = this.get('generateWithPGP');
    const nonce = this.get('nonce');

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
      otp: data.otp,
    };
  },

  attemptProgress(data) {
    const checkStatus = data ? false : true;
    let action = this.get('action');
    action = action && camelize(action);
    this.set('loading', true);
    const adapter = this.get('store').adapterFor('cluster');
    const method = adapter[action];
    method
      .call(adapter, data, { checkStatus })
      .then(resp => this.actionSuccess(resp), (...args) => this.actionError(...args));
  },

  actions: {
    onSubmit(data) {
      if (!data.key) {
        return;
      }
      this.attemptProgress(this.extractData(data));
    },

    startGenerate(data) {
      this.attemptProgress(this.extractData(data));
    },

    generateOTP() {
      const bytes = new window.Uint8Array(16);
      window.crypto.getRandomValues(bytes);
      this.set('otp', base64js.fromByteArray(bytes));
    },

    setKey(_, keyFile) {
      this.set('pgp_key', keyFile.value);
      this.set('pgpKeyFile', keyFile);
    },

    clearToken() {
      this.set('encoded_token', null);
    },
  },
});
