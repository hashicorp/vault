import Ember from 'ember';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { task } from 'ember-concurrency';
const BACKENDS = supportedAuthBackends();
const { computed, inject, get } = Ember;

const DEFAULTS = {
  token: null,
  username: null,
  password: null,
};

export default Ember.Component.extend(DEFAULTS, {
  classNames: ['auth-form'],
  router: inject.service(),
  auth: inject.service(),
  flashMessages: inject.service(),
  store: inject.service(),
  csp: inject.service('csp-event'),

  // set during init and potentially passed in via a query param
  selectedAuth: null,
  methods: null,
  cluster: null,
  redirectTo: null,

  didRender() {
    this._super(...arguments);
    // on very narrow viewports the active tab may be overflowed, so we scroll it into view here
    let activeEle = this.element.querySelector('li.is-active');
    if (activeEle) {
      activeEle.scrollIntoView();
    }
    activeEle = null;
    // this is here because we're changing the `with` attr and there's no way to short-circuit rendering,
    // so we'll just nav -> get new attrs -> re-render
    if (!this.get('selectedAuth') || (this.get('selectedAuth') && !this.get('selectedAuthBackend'))) {
      this.get('router').replaceWith('vault.cluster.auth', this.get('cluster.name'), {
        queryParams: {
          with: this.firstMethod(),
          wrappedToken: this.get('wrappedToken'),
        },
      });
    }
  },

  firstMethod() {
    let firstMethod = this.get('methodsToShow.firstObject');
    // prefer backends with a path over those with a type
    return get(firstMethod, 'path') || get(firstMethod, 'type');
  },

  didReceiveAttrs() {
    this._super(...arguments);
    let token = this.get('wrappedToken');
    let newMethod = this.get('selectedAuth');
    let oldMethod = this.get('oldSelectedAuth');

    if (oldMethod && oldMethod !== newMethod) {
      this.resetDefaults();
    }
    this.set('oldSelectedAuth', newMethod);

    if (token) {
      this.get('unwrapToken').perform(token);
    }
  },

  resetDefaults() {
    this.setProperties(DEFAULTS);
  },

  selectedAuthIsPath: computed.match('selectedAuth', /\/$/),
  selectedAuthBackend: Ember.computed(
    'allSupportedMethods',
    'selectedAuth',
    'selectedAuthIsPath',
    function() {
      let methods = this.get('allSupportedMethods');
      let keyIsPath = this.get('selectedAuthIsPath');
      let findKey = keyIsPath ? 'path' : 'type';
      return methods.findBy(findKey, this.get('selectedAuth'));
    }
  ),

  providerPartialName: computed('selectedAuthBackend', function() {
    let type = this.get('selectedAuthBackend.type') || 'token';
    type = type.toLowerCase();
    let templateName = Ember.String.dasherize(type);
    return `partials/auth-form/${templateName}`;
  }),

  hasCSPError: computed.alias('csp.connectionViolations.firstObject'),

  cspErrorText: `This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`,

  allSupportedMethods: computed('methodsToShow', 'hasMethodsWithPath', function() {
    let hasMethodsWithPath = this.get('hasMethodsWithPath');
    let methodsToShow = this.get('methodsToShow');
    return hasMethodsWithPath ? methodsToShow.concat(BACKENDS) : methodsToShow;
  }),

  hasMethodsWithPath: computed('methodsToShow', function() {
    return this.get('methodsToShow').isAny('path');
  }),
  methodsToShow: computed('methods', 'methods.[]', function() {
    let methods = this.get('methods') || [];
    let shownMethods = methods.filter(m =>
      BACKENDS.find(b => get(b, 'type').toLowerCase() === get(m, 'type').toLowerCase())
    );
    return shownMethods.length ? shownMethods : BACKENDS;
  }),

  unwrapToken: task(function*(token) {
    // will be using the token auth method, so set it here
    this.set('selectedAuth', 'token');
    let adapter = this.get('store').adapterFor('tools');
    try {
      let response = yield adapter.toolAction('unwrap', null, { clientToken: token });
      this.set('token', response.auth.client_token);
      this.send('doSubmit');
    } catch (e) {
      this.set('error', `Token unwrap failed: ${e.errors[0]}`);
    }
  }),

  handleError(e) {
    this.set('loading', false);
    let errors = e.errors.map(error => {
      if (error.detail) {
        return error.detail;
      }
      return error;
    });
    this.set('error', `Authentication failed: ${errors.join('.')}`);
  },

  actions: {
    doSubmit() {
      let data = {};
      this.setProperties({
        loading: true,
        error: null,
      });
      let targetRoute = this.get('redirectTo') || 'vault.cluster';
      let backend = this.get('selectedAuthBackend');
      let path = get(backend, 'path') || this.get('customPath');
      let backendMeta = BACKENDS.find(
        b => get(b, 'type').toLowerCase() === get(backend, 'type').toLowerCase()
      );
      let attributes = get(backendMeta, 'formAttributes');

      data = Ember.assign(data, this.getProperties(...attributes));
      if (get(backend, 'path') || (this.get('useCustomPath') && path)) {
        data.path = path;
      }
      const clusterId = this.get('cluster.id');
      this.get('auth').authenticate({ clusterId, backend: get(backend, 'type'), data }).then(
        ({ isRoot }) => {
          this.set('loading', false);
          const transition = this.get('router').transitionTo(targetRoute);
          if (isRoot) {
            transition.followRedirects().then(() => {
              this.get('flashMessages').warning(
                'You have logged in with a root token. As a security precaution, this root token will not be stored by your browser and you will need to re-authenticate after the window is closed or refreshed.'
              );
            });
          }
        },
        (...errArgs) => this.handleError(...errArgs)
      );
    },
  },
});
