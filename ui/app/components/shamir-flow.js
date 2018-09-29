import { inject as service } from '@ember/service';
import { gt } from '@ember/object/computed';
import { camelize } from '@ember/string';
import Component from '@ember/component';
import { get, computed } from '@ember/object';
import base64js from 'base64-js';

const DEFAULTS = {
  key: null,
  loading: false,
  errors: [],
  threshold: null,
  progress: null,
  pgp_key: null,
  haveSavedPGPKey: false,
  started: false,
  generateWithPGP: false,
  pgpKeyFile: { value: '' },
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

  init() {
    this._super(...arguments);
    if (this.get('fetchOnInit')) {
      this.attemptProgress();
    }
  },

  didInsertElement() {
    this._super(...arguments);
    this.onUpdate(this.getProperties(Object.keys(DEFAULTS)));
  },

  onUpdate() {},
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
    let { onUpdate, isComplete, onShamirSuccess, thresholdPath } = this.getProperties(
      'onUpdate',
      'isComplete',
      'onShamirSuccess',
      'thresholdPath'
    );
    let threshold = get(resp, thresholdPath);
    let props = {
      ...resp,
      threshold,
    };
    this.stopLoading();
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
      throw e;
    }
  },

  generateStep: computed('generateWithPGP', 'haveSavedPGPKey', 'otp', 'pgp_key', function() {
    let { generateWithPGP, otp, pgp_key, haveSavedPGPKey } = this.getProperties(
      'generateWithPGP',
      'otp',
      'pgp_key',
      'haveSavedPGPKey'
    );
    if (!generateWithPGP && !pgp_key && !otp) {
      return 'chooseMethod';
    }
    if (otp) {
      return 'beginGenerationWithOTP';
    }
    if (generateWithPGP) {
      if (pgp_key && haveSavedPGPKey) {
        return 'beginGenerationWithPGP';
      } else {
        return 'providePGPKey';
      }
    }
  }),

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

    savePGPKey() {
      if (this.get('pgp_key')) {
        this.set('haveSavedPGPKey', true);
      }
    },
  },
});
