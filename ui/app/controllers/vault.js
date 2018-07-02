import Ember from 'ember';

export default Ember.Controller.extend({
  queryParams: [
    {
      wrappedToken: 'wrapped_token',
    },
  ],
  wrappedToken: '',
});
