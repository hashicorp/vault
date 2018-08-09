import Ember from 'ember';

export default Ember.Controller.extend({
  vaultController: Ember.inject.controller('vault'),
  queryParams: [{ authMethod: 'with' }],
  wrappedToken: Ember.computed.alias('vaultController.wrappedToken'),
  authMethod: '',
  redirectTo: null,
});
