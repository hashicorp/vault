import Ember from 'ember';
import { task } from 'ember-concurrency';

const { computed, inject } = Ember;

export default Ember.Component.extend({
  tagName: '',
  auth: inject.service(),
  controlGroup: inject.service(),

  // public API
  model: null,

  init() {
    this._super(...arguments);
    let accessor = this.get('model.id');
    let data = this.get('controlGroup').wrapInfoForAccessor(accessor);
    if (data) {
      this.set('controlGroupResponse', data);
    }
  },

  currentUserEntityId: computed.alias('auth.authData.entity_id'),

  currentUserIsRequesting: computed('currentUserEntityId', 'model.requestEntity.id', function() {
    return this.get('currentUserEntityId') === this.get('model.requestEntity.id');
  }),

  currentUserHasAuthorized: computed('currentUserEntityId', 'model.authorizations.@each.id', function() {
    let authorizations = this.get('model.authorizations') || [];
    return Boolean(authorizations.findBy('id', this.get('currentUserEntityId')));
  }),

  isSuccess: computed.or('currentUserHasAuthorized', 'model.approved'),
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

  bannerPrefix: computed('model.approved', 'currentUserIsRequesting', 'currentUserHasAuthorized', function() {
    let isApproved = this.get('model.approved');
    let { currentUserHasAuthorized, currentUserIsRequesting } = this.getProperties(
      'currentUserIsRequesting',
      'currentUserHasAuthorized'
    );
    if (currentUserHasAuthorized) {
      return 'Thanks!';
    }
    if (isApproved) {
      return 'Success!';
    }
    return 'Locked';
  }),

  bannerText: computed('model.approved', 'currentUserIsRequesting', 'currentUserHasAuthorized', function() {
    let isApproved = this.get('model.approved');
    let { currentUserHasAuthorized, currentUserIsRequesting } = this.getProperties(
      'currentUserIsRequesting',
      'currentUserHasAuthorized'
    );
    if (currentUserHasAuthorized) {
      return 'You have given authorization';
    }
    if (currentUserIsRequesting && isApproved) {
      return 'You have been given authorization';
    }
    if (isApproved) {
      return 'This control group has been authorized';
    }
    if (currentUserIsRequesting) {
      return 'The path you requested is locked by a control group';
    }
    return 'Someone is requesting access to a path locked by a control group';
  }),

  refresh: task(function*() {
    yield this.get('model').reload();
  }).drop(),

  authorize: task(function*() {
    try {
      yield this.get('model').save();
      yield this.get('refresh').perform();
    } catch (e) {
      this.set('errors', e);
    }
  }).drop(),
});
