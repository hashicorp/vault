import Ember from 'ember';
import moment from 'moment';

const { get, set, computed, setProperties } = Ember;

const DEFAULTS = {
  token: null,
  rewrap_token: null,
  errors: [],
  wrap_info: null,
  creation_time: null,
  creation_ttl: null,
  data: '{\n}',
  unwrap_data: null,
  wrapTTL: null,
  sum: null,
  random_bytes: null,
  input: null,
};

const WRAPPING_ENDPOINTS = ['lookup', 'wrap', 'unwrap', 'rewrap'];

export default Ember.Component.extend(DEFAULTS, {
  store: Ember.inject.service(),
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
    const currentAction = get(this, 'selectedAction');
    const oldAction = get(this, 'oldSelectedAction');

    if (currentAction !== oldAction) {
      this.reset();
    }
    set(this, 'oldSelectedAction', currentAction);
  },

  dataIsEmpty: computed.match('data', new RegExp(DEFAULTS.data)),

  expirationDate: computed('creation_time', 'creation_ttl', function() {
    const { creation_time, creation_ttl } = this.getProperties('creation_time', 'creation_ttl');
    if (!(creation_time && creation_ttl)) {
      return null;
    }
    return moment(creation_time).add(moment.duration(creation_ttl, 'seconds'));
  }),

  handleError(e) {
    set(this, 'errors', e.errors);
  },

  handleSuccess(resp, action) {
    let props = {};
    let secret = (resp && resp.data) || resp.auth;
    if (secret && action === 'unwrap') {
      props = Ember.assign({}, props, { unwrap_data: secret });
    }
    props = Ember.assign({}, props, secret);

    if (resp && resp.wrap_info) {
      const keyName = action === 'rewrap' ? 'rewrap_token' : 'token';
      props = Ember.assign({}, props, { [keyName]: resp.wrap_info.token });
    }
    setProperties(this, props);
  },

  getData() {
    const action = get(this, 'selectedAction');
    if (WRAPPING_ENDPOINTS.includes(action)) {
      return get(this, 'dataIsEmpty')
        ? { token: (get(this, 'token') || '').trim() }
        : JSON.parse(get(this, 'data'));
    }
    if (action === 'random') {
      return this.getProperties('bytes');
    }

    if (action === 'hash') {
      return this.getProperties('input', 'format', 'algorithm');
    }
  },

  actions: {
    doSubmit(evt) {
      evt.preventDefault();
      const action = get(this, 'selectedAction');
      const wrapTTL = action === 'wrap' ? get(this, 'wrapTTL') : null;
      const data = this.getData();
      setProperties(this, {
        errors: null,
        wrap_info: null,
        creation_time: null,
        creation_ttl: null,
      });

      get(this, 'store')
        .adapterFor('tools')
        .toolAction(action, data, { wrapTTL })
        .then(resp => this.handleSuccess(resp, action), (...errArgs) => this.handleError(...errArgs));
    },

    onClear() {
      this.reset();
    },

    codemirrorUpdated(val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;
      setProperties(this, {
        buttonDisabled: hasErrors,
        data: val,
      });
    },
  },
});
