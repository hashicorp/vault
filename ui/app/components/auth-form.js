import Ember from 'ember';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
const BACKENDS = supportedAuthBackends();
const { computed, inject, get } = Ember;

const DEFAULTS = {
  token: null,
  username: null,
  password: null,
};

export default Ember.Component.extend(DEFAULTS, {
  classNames: ['auth-form'],
  routing: inject.service('-routing'),
  auth: inject.service(),
  flashMessages: inject.service(),
  csp: inject.service('csp-event'),
  didRender() {
    // on very narrow viewports the active tab may be overflowed, so we scroll it into view here
    this.$('li.is-active').get(0).scrollIntoView();
  },

  didReceiveAttrs() {
    this._super(...arguments);
    let newMethod = this.get('selectedAuthType');
    let oldMethod = this.get('oldSelectedAuthType');

    if (oldMethod && oldMethod !== newMethod) {
      this.resetDefaults();
    }
    this.set('oldSelectedAuthType', newMethod);
  },

  resetDefaults() {
    this.setProperties(DEFAULTS);
  },

  cluster: null,
  redirectTo: null,

  selectedAuthType: 'token',
  selectedAuthBackend: Ember.computed('selectedAuthType', function() {
    return BACKENDS.findBy('type', this.get('selectedAuthType'));
  }),

  providerPartialName: Ember.computed('selectedAuthType', function() {
    const type = Ember.String.dasherize(this.get('selectedAuthType'));
    return `partials/auth-form/${type}`;
  }),

  hasCSPError: computed.alias('csp.connectionViolations.firstObject'),

  cspErrorText: `This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`,

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
      let path = this.get('customPath');
      let attributes = get(backend, 'formAttributes');

      data = Ember.assign(data, this.getProperties(...attributes));
      if (this.get('useCustomPath') && path) {
        data.path = path;
      }
      const clusterId = this.get('cluster.id');
      this.get('auth').authenticate({ clusterId, backend: get(backend, 'type'), data }).then(
        ({ isRoot }) => {
          this.set('loading', false);
          const transition = this.get('routing.router').transitionTo(targetRoute);
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
