import Ember from 'ember';

const { inject, Controller } = Ember;
export default Controller.extend({
  namespaceService: inject.service('namespace'),
});
