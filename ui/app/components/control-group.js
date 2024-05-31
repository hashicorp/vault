/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { alias, or } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { task } from 'ember-concurrency';

export default Component.extend({
  tagName: '',
  auth: service(),
  controlGroup: service(),

  // public API
  model: null,

  didReceiveAttrs() {
    this._super(...arguments);
    const accessor = this.model.id;
    const data = this.controlGroup.wrapInfoForAccessor(accessor);
    this.set('controlGroupResponse', data);
  },

  currentUserEntityId: alias('auth.authData.entity_id'),

  currentUserIsRequesting: computed('currentUserEntityId', 'model.requestEntity.id', function () {
    if (!this.model.requestEntity) return false;
    return this.currentUserEntityId === this.model.requestEntity.id;
  }),

  currentUserHasAuthorized: computed('currentUserEntityId', 'model.authorizations.@each.id', function () {
    const authorizations = this.model.authorizations || [];
    return Boolean(authorizations.find((authz) => authz.id === this.currentUserEntityId));
  }),

  isSuccess: or('currentUserHasAuthorized', 'model.approved'),
  requestorName: computed('currentUserIsRequesting', 'model.requestEntity', function () {
    const entity = this.model.requestEntity;

    if (this.currentUserIsRequesting) {
      return 'You';
    }
    if (entity && entity.name) {
      return entity.name;
    }
    return 'Someone';
  }),

  bannerPrefix: computed('model.approved', 'currentUserHasAuthorized', function () {
    if (this.currentUserHasAuthorized) {
      return 'Thanks!';
    }
    if (this.model.approved) {
      return 'Success!';
    }
    return 'Locked';
  }),

  bannerText: computed('model.approved', 'currentUserIsRequesting', 'currentUserHasAuthorized', function () {
    const isApproved = this.model.approved;
    const { currentUserHasAuthorized, currentUserIsRequesting } = this;
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

  refresh: task(function* () {
    try {
      yield this.model.reload();
    } catch (e) {
      this.set('errors', e);
    }
  }).drop(),

  authorize: task(function* () {
    try {
      yield this.model.save();
      yield this.refresh.perform();
    } catch (e) {
      this.set('errors', e);
    }
  }).drop(),
});
