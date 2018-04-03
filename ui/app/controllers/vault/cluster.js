import Ember from 'ember';

const { inject, Controller } = Ember;

export default Controller.extend({
  auth: inject.service(),
  version: inject.service(),
});
