import Ember from 'ember';

export default Ember.Test.registerAsyncHelper('authLogout', function() {
  visit('/vault/logout');
});
