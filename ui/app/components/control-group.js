import { inject as service } from '@ember/service';
import { alias, or } from '@ember/object/computed';
import Component from '@ember/component';
import { computed, get } from '@ember/object';
import { task } from 'ember-concurrency';

export default Component.extend({
  tagName: '',
  auth: service(),
  controlGroup: service(),

  // public API
  model: null,

  didReceiveAttrs() {
    this._super(...arguments);
    let accessor = this.get('model.id');
    let data = this.get('controlGroup').wrapInfoForAccessor(accessor);
    this.set('controlGroupResponse', data);
  },

  currentUserEntityId: alias('auth.authData.entity_id'),

  currentUserIsRequesting: computed('currentUserEntityId', 'model.requestEntity.id', function() {
    return this.get('currentUserEntityId') === this.get('model.requestEntity.id');
  }),

  currentUserHasAuthorized: computed('currentUserEntityId', 'model.authorizations.@each.id', function() {
    let authorizations = this.get('model.authorizations') || [];
    return Boolean(authorizations.findBy('id', this.get('currentUserEntityId')));
  }),

  isSuccess: or('currentUserHasAuthorized', 'model.approved'),
  requestorName: computed('currentUserIsRequesting', 'model.requestEntity', function() {
    let entity = this.get('model.requestEntity');

    if (this.get('currentUserIsRequesting')) {
      return 'You';
    }
    if (entity && get(entity, 'name')) {
      return get(entity, 'name');
    }
    return 'Someone';
  }),

  bannerPrefix: computed('model.approved', 'currentUserHasAuthorized', function() {
    if (this.get('currentUserHasAuthorized')) {
      return 'Thanks!';
    }
    if (this.get('model.approved')) {
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
      return 'This Control Group has been authorized';
    }
    if (currentUserIsRequesting) {
      return 'The path you requested is locked by a Control Group';
    }
    return 'Someone is requesting access to a path locked by a Control Group';
  }),

  refresh: task(function*() {
    try {
      yield this.get('model').reload();
    } catch (e) {
      this.set('errors', e);
    }
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
