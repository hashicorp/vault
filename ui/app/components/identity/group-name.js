/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

export default class IdentityGroupNameComponent extends Component {
  @service store;

  @tracked groupName = null;
  @tracked isLoading = true;
  @tracked hasError = false;

  constructor() {
    super(...arguments);
    this.loadGroupName.perform();
  }

  @task
  *loadGroupName() {
    try {
      this.isLoading = true;
      this.hasError = false;
      
      // First try to find the record in the store cache
      let group = this.store.peekRecord('identity/group', this.args.groupId);
      
      if (!group) {
        // If not in cache, fetch from API
        group = yield this.store.findRecord('identity/group', this.args.groupId);
      }
      
      this.groupName = group.name;
    } catch (error) {
      this.hasError = true;
      // Fallback to showing the ID if we can't load the group
      this.groupName = this.args.groupId;
    } finally {
      this.isLoading = false;
    }
  }

  get displayName() {
    if (this.isLoading) {
      return 'Loading...';
    }
    return this.groupName || this.args.groupId;
  }
}