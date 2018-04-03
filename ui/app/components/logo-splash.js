import Ember from 'ember';

const { inject } = Ember;

export default Ember.Component.extend({
  tagName: '',
  version: inject.service(),
});
