import Ember from 'ember';

const { inject } = Ember;

export default Ember.Component.extend({
  version: inject.service(),
  tagName: '',
});
