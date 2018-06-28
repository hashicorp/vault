import Ember from 'ember';

export default Ember.Controller.extend({
  queryParams: [
    'with',
    {
      wrappedToken: 'wrapped_token',
    },
  ],
  wrappedToken: '',
  with: '',
  redirectTo: null,
});
