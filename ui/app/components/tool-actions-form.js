import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action, setProperties } from '@ember/object';
import { assign } from '@ember/polyfills';
/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import { match } from '@ember/object/computed';
import { tracked } from '@glimmer/tracking';
import { addSeconds, parseISO } from 'date-fns';
import { A } from '@ember/array';

/**
 * @module ToolActionsForm2
 * ToolActionsForm2 components are used to...
 *
 * @example
 * ```js
 * <ToolActionsForm2 @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export const DEFAULTS = {
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

export const WRAPPING_ENDPOINTS = ['lookup', 'wrap', 'unwrap', 'rewrap'];

export default class ToolActionFormTwo extends Component {
  @service store;
  @service wizard;

  @tracked bytes = 32;
  @tracked format = 'base64';
  @tracked alogrithm = 'sha2-256';
  @tracked unwrapActiveTab = 'data';
  @tracked creation_time;
  @tracked creation_ttl;

  constructor() {
    super(...arguments);
    this.checkAction();
  }

  reset() {
    if (this.isDestroyed || this.isDestroying) {
      return;
    }
    this.DEFAULTS = setProperties(this, DEFAULTS);
  }

  checkAction() {
    if (this.args.selectedAction !== this.oldSelectedAction) {
      this.reset();
    }
    this.oldSelectedAction = this.args.selectedAction;
  }
  @match('data', new RegExp(DEFAULTS.data)) dataIsEmpty;

  get expirationDate() {
    const { creation_time, creation_ttl } = this;
    if (!(creation_time && creation_ttl)) {
      return null;
    }
    // returns new Date with seconds added.
    return addSeconds(parseISO(creation_time), creation_ttl);
  }

  handleError(e) {
    this.errors = e.errors;
  }

  handleSuccess(resp, action) {
    let props = {};
    let secret = (resp && resp.data) || resp.auth;
    if (secret && action === 'unwrap') {
      let details = {
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
    if (props.token || props.rewrap_token || props.unwrap_data || action === 'lookup') {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE');
    }
    setProperties(this, props);
  }

  getData() {
    const action = this.args.selectedAction;
    if (WRAPPING_ENDPOINTS.includes(action)) {
      return this.dataIsEmpty ? { token: (this.token || '').trim() } : JSON.parse(this.data);
    }
    if (action === 'random') {
      return { bytes: this.bytes, format: this.format };
    }
    if (action === 'hash') {
      return { input: this.input, format: this.format, algorithm: this.algorithm };
    }
  }

  @action
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
  }

  @action
  onClear() {
    this.reset();
  }

  @action
  onHash(newValue) {
    this.input = newValue;
  }

  @action
  updateTtl(ttl) {
    this.wrapTTl = ttl;
  }

  @action
  codemirrorUpdated(val, hasErrors) {
    setProperties(this, {
      buttonDisabled: hasErrors,
      data: val,
    });
  }
}
