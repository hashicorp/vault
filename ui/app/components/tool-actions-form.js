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
 * @module ToolActionsForm
 * ToolActionsForm components is used to display the tool fool under to tools/tool route.
 *
 * @example
 * ```js
 * <ToolActionsForm @selectedAction=this.selectedAction/>
 * ```
 * @param {string} selectedAction - The action selected such as wrap, lookup, etc.

 */

export const DEFAULTS = {
  creation_time: null,
  creation_ttl: null,
  data: '{\n}',
  details: null,
  errors: A(),
  input: null,
  random_bytes: null,
  rewrap_token: null,
  sum: null,
  token: null,
  unwrap_data: null,
  wrap_info: null,
};

export const WRAPPING_ENDPOINTS = ['lookup', 'wrap', 'unwrap', 'rewrap'];

export default class ToolActionForm extends Component {
  @service store;
  @service wizard;

  @tracked algorithm = 'sha2-256';
  @tracked bytes = 32;
  @tracked creation_time;
  @tracked creation_ttl;
  @tracked errors;
  @tracked format = 'base64';
  @tracked oldSelectedAction;
  @tracked input;
  @tracked unwrapActiveTab = 'data';
  @tracked wrapTTL = '30m'; // default unless otherwise selected

  @match('data', new RegExp(DEFAULTS.data)) dataIsEmpty;

  constructor() {
    super(...arguments);
    this.checkAction();
  }

  get expirationDate() {
    const { creation_time, creation_ttl } = this;
    if (!(creation_time && creation_ttl)) {
      return null;
    }
    // returns new Date with seconds added.
    return addSeconds(parseISO(creation_time), creation_ttl);
  }

  checkAction() {
    if (this.args.selectedAction !== this.oldSelectedAction) {
      this.reset();
    }
    this.oldSelectedAction = this.args.selectedAction;
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

  reset() {
    if (this.isDestroyed || this.isDestroying) {
      return;
    }
    this.DEFAULTS = setProperties(this, DEFAULTS);
  }

  get toolAdapter() {
    return this.store.adapterFor('tools');
  }

  @action
  codemirrorUpdated(val, hasErrors) {
    setProperties(this, {
      buttonDisabled: hasErrors,
      data: val,
    });
  }

  @action
  async doSubmit(evt) {
    evt.preventDefault();
    const action = this.args.selectedAction;
    const wrapTTL = action === 'wrap' ? this.wrapTTL : null;
    const data = this.getData();
    setProperties(this, {
      errors: null,
      wrap_info: null,
      creation_time: null,
      creation_ttl: null,
    });

    try {
      let response = await this.toolAdapter.toolAction(action, data, { wrapTTL });
      this.handleSuccess(response, action);
    } catch (error) {
      this.handleError(error);
    }
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
}
