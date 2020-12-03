import { next } from '@ember/runloop';
import { inject as service } from '@ember/service';
import { match, alias, or } from '@ember/object/computed';
import { assign } from '@ember/polyfills';
import { dasherize } from '@ember/string';
import Component from '@ember/component';
import { get, computed } from '@ember/object';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { task } from 'ember-concurrency';
const BACKENDS = supportedAuthBackends();

import { toolsActions } from 'vault/helpers/tools-actions';

/**
 * @module AuthForm
 * The `AuthForm` is used to sign users into Vault.
 *
 * @example ```js
 * // All properties are passed in via query params.
 *   <AuthForm @wrappedToken={{wrappedToken}} @cluster={{model}} @namespace={{namespaceQueryParam}} @redirectTo={{redirectTo}} @selectedAuth={{authMethod}}/>```
 *
 * @param wrappedToken=null {String} - The auth method that is currently selected in the dropdown.
 * @param cluster=null {Object} - The auth method that is currently selected in the dropdown. This corresponds to an Ember Model.
 * @param namespace=null {String} - The currently active namespace.
 * @param redirectTo=null {String} - The name of the route to redirect to.
 * @param selectedAuth=null {String} - The auth method that is currently selected in the dropdown.
 */

// const DEFAULTS = {
//   token: null,
//   username: null,
//   password: null,
//   customPath: null,
// };
const DEFAULTS = {
  token: null,
  rewrap_token: null,
  errors: [],
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
  username: null,
  password: null,
  customPath: null,
};
const WRAPPING_ENDPOINTS = ['lookup', 'wrap', 'unwrap', 'rewrap'];

export default Component.extend(DEFAULTS, {
  // putting these attrs here so they don't get reset when you click back
  //random
  bytes: 32,
  //hash
  format: 'base64',
  algorithm: 'sha2-256',
  tagName: '',
  unwrapActiveTab: 'data',
  selectedAction: null,

  router: service(),
  auth: service(),
  flashMessages: service(),
  store: service(),
  csp: service('csp-event'),
  currentCluster: service(),
  wizard: service(),

  //  passed in via a query param
  selectedAuth: null,
  methods: null,
  cluster: null,
  redirectTo: null,
  namespace: null,
  wrappedToken: null,
  // internal
  oldNamespace: null,

  beforeModel(transition) {
    const currentCluster = this.get('currentCluster.cluster.name');
    const { selected_action: selectedAction } = this.paramsFor(this.routeName);
    const supportedActions = toolsActions();

    if (transition.targetName === this.routeName) {
      transition.abort();
      return this.replaceWith('vault.cluster.tools.tool', currentCluster, supportedActions[2]);
    }
  },
  model(params) {
    return params.selected_action;
  },

  setupController(controller, model) {
    this._super(...arguments);
    controller.set('selectedAction', model);
  },

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
      this.get('wizard').transitionFeatureMachine(this.get('wizard.featureState'), 'CONTINUE');
    }
    this.setProperties(props);
  },

  onClear() {
    this.reset();
  },

  updateTtl(evt) {
    const ttl = evt.enabled ? `${evt.seconds}s` : '30m';
    set(this, 'wrapTTL', ttl);
  },

  codemirrorUpdated(val, codemirror) {
    codemirror.performLint();
    const hasErrors = codemirror.state.lint.marked.length > 0;
    this.setProperties(this, {
      buttonDisabled: hasErrors,
      data: val,
    });
  },

  reset() {
    if (this.selectedAuth !== 'unwrap') return ;
    if (this.isDestroyed || this.isDestroying) {
      return;
    }
    this.setProperties(this, DEFAULTS);
  },

  checkAction() {
    const currentAction = get(this, 'selectedAction');
    const oldAction = get(this, 'oldSelectedAction');
    if (currentAction === undefined) {
      return;
    }
    if (currentAction !== oldAction) {
      this.reset();
    }
    this.set('oldSelectedAction', currentAction);
  },

  didReceiveAttrs() {
    let {
      wrappedToken: token,
      oldWrappedToken: oldToken,
      oldNamespace: oldNS,
      namespace: ns,
      selectedAuth: newMethod,
      oldSelectedAuth: oldMethod,
    } = this;

    this._super(...arguments);
    if (newMethod === "unwrap" && oldMethod !== undefined ) this.checkAction();
    else
    {
      next(() => {
        if (!token && (oldNS === null || oldNS !== ns)) {
          this.fetchMethods.perform();
        }
        if (ns !== undefined && ns !== "") this.set('oldNamespace', ns);

        // we only want to trigger this once
        if (token && !oldToken) {
          this.unwrapToken.perform(token);
          this.set('oldWrappedToken', token);
        }
        if (oldMethod && oldMethod !== newMethod && newMethod !== 'unwrap') {
          this.resetDefaults();
        }
        if (newMethod !== 'unwrap') this.set('oldSelectedAuth', newMethod);
      });
    }
  },

  didRender() {
    if (this.element === null) return ;
    //if (this.authMethod === 'unwrap' || this.authMethod === undefined) return ;
    this._super(...arguments);
    // on very narrow viewports the active tab may be overflowed, so we scroll it into view here
    let activeEle = this.element.querySelector('li.is-active');
    if (activeEle) {
      activeEle.scrollIntoView();
    }

    next(() => {
      let firstMethod = this.firstMethod();
      // set `with` to the first method
      if (
        !this.wrappedToken &&
        ((this.get('fetchMethods.isIdle') && firstMethod && !this.get('selectedAuth')) ||
          (this.get('selectedAuth') && !this.get('selectedAuthBackend')))
      ) {
        this.set('selectedAuth', firstMethod);
      }
    });
  },

  firstMethod() {
    let firstMethod = this.get('methodsToShow.firstObject');
    if (!firstMethod) return;
    // prefer backends with a path over those with a type
    return get(firstMethod, 'path') || get(firstMethod, 'type');
  },

  resetDefaults() {
    this.setProperties(DEFAULTS);
  },

  selectedAuthIsPath: match('selectedAuth', /\/$/),
  selectedAuthBackend: computed(
    'wrappedToken',
    'methods',
    'methods.[]',
    'selectedAuth',
    'selectedAuthIsPath',
    function() {
      let { wrappedToken, methods, selectedAuth, selectedAuthIsPath: keyIsPath } = this;
      if (!methods && !wrappedToken) {
        return {};
      }
      if (keyIsPath) {
        return methods.findBy('path', selectedAuth);
      }
      return BACKENDS.findBy('type', selectedAuth);
    }
  ),

  providerPartialName: computed('selectedAuthBackend', function() {
    let type = this.get('selectedAuthBackend.type') || 'token';
    type = type.toLowerCase();
    let templateName = dasherize(type);
    return `partials/auth-form/${templateName}`;
  }),

  hasCSPError: alias('csp.connectionViolations.firstObject'),

  cspErrorText: `This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`,

  allSupportedMethods: computed('methodsToShow', 'hasMethodsWithPath', function() {
    let hasMethodsWithPath = this.get('hasMethodsWithPath');
    let methodsToShow = this.get('methodsToShow');
    return hasMethodsWithPath ? methodsToShow.concat(BACKENDS) : methodsToShow;
  }),

  hasMethodsWithPath: computed('methodsToShow', function() {
    return this.get('methodsToShow').isAny('path');
  }),
  methodsToShow: computed('methods', function() {
    let methods = this.get('methods') || [];
    let shownMethods = methods.filter(m =>
      BACKENDS.find(b => get(b, 'type').toLowerCase() === get(m, 'type').toLowerCase())
    );
    return shownMethods.length ? shownMethods : BACKENDS;
  }),

  unwrapToken: task(function*(token) {
    // will be using the Token Auth Method, so set it here
    this.set('selectedAuth', 'token');
    let adapter = this.get('store').adapterFor('tools');
    try {
      let response = yield adapter.toolAction('unwrap', null, { clientToken: token });
      this.set('token', response.auth.client_token);
      this.send('doSubmit');
    } catch (e) {
      this.set('error', `Token unwrap failed: ${e.errors[0]}`);
    }
  }).withTestWaiter(),

  fetchMethods: task(function*() {
    let store = this.get('store');
    try {
      let methods = yield store.findAll('auth-method', {
        adapterOptions: {
          unauthenticated: true,
        },
      });
      this.set('methods', methods.map(m => m.serialize({ includeId: true })));
      next(() => {
        store.unloadAll('auth-method');
      });
    } catch (e) {
      this.set('error', `There was an error fetching Auth Methods: ${e.errors[0]}`);
    }
  }).withTestWaiter(),

  showLoading: or('isLoading', 'authenticate.isRunning', 'fetchMethods.isRunning', 'unwrapToken.isRunning'),

  handleError(e, prefixMessage = true) {
    if (this.authMethod === undefined){
      this.set('errors', e.errors);
      return ;
    }

    this.set('loading', false);
    let errors;
    if (e.errors) {
      errors = e.errors.map(error => {
        if (error.detail) {
          return error.detail;
        }
        return error;
      });
    } else {
      errors = [e];
    }
    let message = prefixMessage ? 'Authentication failed: ' : '';
    this.set('error', `${message}${errors.join('.')}`);
  },

  authenticate: task(function*(backendType, data) {
    let clusterId = this.cluster.id;
    try {
      let authResponse = yield this.auth.authenticate({ clusterId, backend: backendType, data });
      let { isRoot, namespace } = authResponse;
      let transition;
      let { redirectTo } = this;
      if (redirectTo) {
        // reset the value on the controller because it's bound here
        this.set('redirectTo', '');
        // here we don't need the namespace because it will be encoded in redirectTo
        transition = this.router.transitionTo(redirectTo);
      } else {
        transition = this.router.transitionTo('vault.cluster', { queryParams: { namespace } });
      }
      // returning this w/then because if we keep it
      // in the task, it will get cancelled when the component in un-rendered
      yield transition.followRedirects().then(() => {
        if (isRoot) {
          this.flashMessages.warning(
            'You have logged in with a root token. As a security precaution, this root token will not be stored by your browser and you will need to re-authenticate after the window is closed or refreshed.'
          );
        }
      });
    } catch (e) {
      this.handleError(e);
    }
  }).withTestWaiter(),

  actions: {
    didTransition() {
      const params = this.paramsFor(this.routeName);
      if (this.wizard.currentMachine === 'tools') {
        this.wizard.transitionFeatureMachine(this.wizard.featureState, params.selected_action.toUpperCase());
      }
      this.controller.setProperties(params);
      return true;
    },
    doSubmit(evt) {
      if (this.selectedAuth !== "unwrap")
      {
        let passedData, e;
        if (arguments.length > 1) {
          [passedData, e] = arguments;
        } else {
          [e] = arguments;
        }
        if (e) {
          e.preventDefault();
        }
        let data = {};
        this.setProperties({
          error: null,
        });
        let backend = this.get('selectedAuthBackend') || {};
        let backendMeta = BACKENDS.find(
          b => (get(b, 'type') || '').toLowerCase() === (get(backend, 'type') || '').toLowerCase()
        );
        let attributes = get(backendMeta || {}, 'formAttributes') || [];
  
        data = assign(data, this.getProperties(...attributes));
        if (passedData) {
          data = assign(data, passedData);
        }
        if (this.get('customPath') || get(backend, 'id')) {
          data.path = this.get('customPath') || get(backend, 'id');
        }
        return this.authenticate.unlinked().perform(backend.type, data);
      }
      else
      {
        evt.preventDefault();
        const action = this.selectedAction;
        const wrapTTL = action === 'wrap' ? get(this, 'wrapTTL') : null;
        const data = { token: this.get('token') };  // this.getData();
        this.setProperties(this, {
          errors: null,
          wrap_info: null,
          creation_time: null,
          creation_ttl: null,
        });

        this.get('store')
          .adapterFor('tools')
          .toolAction('unwrap', '', { clientToken: data.token })
          .then(resp => this.handleSuccess(resp, 'unwrap'), (...errArgs) => this.handleError(...errArgs));
      }
    },
    onClear() {
      this.reset();
      this.resetDefaults();
    },
    handleError(e) {
      if (e) {
        this.handleError(e, false);
      } else {
        this.set('error', null);
      }
    },
  },
});
