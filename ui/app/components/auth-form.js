import Ember from 'ember';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
const BACKENDS = supportedAuthBackends();
const { inject } = Ember;

export default Ember.Component.extend({
  classNames: ['auth-form'],
  routing: inject.service('-routing'),
  auth: inject.service(),
  flashMessages: inject.service(),
  csp: inject.service('csp-event'),
  didRender() {
    // on very narrow viewports the active tab may be overflowed, so we scroll it into view here
    this.$('li.is-active').get(0).scrollIntoView();
  },

  cluster: null,
  redirectTo: null,

  selectedAuthType: 'token',
  selectedAuthBackend: Ember.computed('selectedAuthType', function() {
    return BACKENDS.findBy('type', this.get('selectedAuthType'));
  }),

  providerComponentName: Ember.computed('selectedAuthBackend.type', function() {
    const type = Ember.String.dasherize(this.get('selectedAuthBackend.type'));
    return `auth-form/${type}`;
  }),

  handleError(e) {
    let cluster = this.get('cluster');
    this.set('loading', false);
    if (this.get('csp.connectionViolations.firstObject') && cluster.get('standby')) {
      // this is a CSP error to a disallowed connect-src domain - having the API redirect will cause this;
      this.set(
        'error',
        `This is a standby Vault node and it appears that request forwarding is
        not properly configured. To use the UI for anything other than unsealing
        this node, you will have to navigate to the active Vault node and authenticate
        there.`
      );
      return;
    }

    let errors = e.errors.map(error => {
      if (error.detail) {
        return error.detail;
      }
      return error;
    });

    this.set('error', `Authentication failed: ${errors.join('.')}`);
  },

  actions: {
    doSubmit(data) {
      this.setProperties({
        loading: true,
        error: null,
      });
      const targetRoute = this.get('redirectTo') || 'vault.cluster';
      //const {password, token, username} = data;
      const backend = this.get('selectedAuthBackend.type');
      const path = this.get('customPath');
      if (this.get('useCustomPath') && path) {
        data.path = path;
      }
      const clusterId = this.get('cluster.id');
      this.get('auth').authenticate({ clusterId, backend, data }).then(
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
