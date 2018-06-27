import Ember from 'ember';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';

export default Ember.Controller.extend({
  queryParams: [
    'with',
    {
      wrappedToken: 'wrapped_token',
    },
  ],
  wrappedToken: '',
  with: Ember.computed(function() {
    return supportedAuthBackends()[0].type;
  }),

  redirectTo: null,
});
