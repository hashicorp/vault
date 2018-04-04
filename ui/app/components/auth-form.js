import Ember from 'ember';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
const BACKENDS = supportedAuthBackends();

export default Ember.Component.extend({
  classNames: ['auth-form'],
  routing: Ember.inject.service('-routing'),
  auth: Ember.inject.service(),
  flashMessages: Ember.inject.service(),
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
    this.set('loading', false);
    this.set('error', `Authentication failed: ${e.errors.join('.')}`);
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
