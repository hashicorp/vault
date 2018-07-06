import Ember from 'ember';
import { task } from 'ember-concurrency';

const { computed, inject } = Ember;

export default Ember.Component.extend({
  tagName: '',
  auth: inject.service(),
  // public API
  model: null,
  onRefresh() {},

  currentUserEntityId: computed.alias('auth.authData.entity_id'),

  currentUserIsRequesting: computed('currentUserEntityId', 'model.requestEntity.id', function() {
    return this.get('currentUserEntityId') === this.get('model.requestEntity.id');
  }),

  requestorName: computed('currentUserIsRequesting', 'model.requestEntity', function() {
    let entity = this.get('model.requestEntity');

    if (this.get('currentUserIsRequesting')) {
      return 'You';
    }
    if (entity && entity.get('name')) {
      return entity.get('name');
    }
    return 'Someone';
  }),

  authorize: task(function*() {
    try {
      yield this.get('model').save();
      this.get('onRefresh')();
    } catch (e) {
      this.set('errors', e);
    }
  }).drop(),
});
