/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

export default class IdentityEntityNameComponent extends Component {
  @service store;

  @tracked entityName = null;
  @tracked isLoading = true;
  @tracked hasError = false;

  constructor() {
    super(...arguments);
    this.loadEntityName.perform();
  }

  @task
  *loadEntityName() {
    try {
      this.isLoading = true;
      this.hasError = false;
      
      // First try to find the record in the store cache
      let entity = this.store.peekRecord('identity/entity', this.args.entityId);
      
      if (!entity) {
        // If not in cache, fetch from API
        entity = yield this.store.findRecord('identity/entity', this.args.entityId);
      }
      
      this.entityName = entity.name;
    } catch (error) {
      this.hasError = true;
      // Fallback to showing the ID if we can't load the entity
      this.entityName = this.args.entityId;
    } finally {
      this.isLoading = false;
    }
  }

  get displayName() {
    if (this.isLoading) {
      return 'Loading...';
    }
    return this.entityName || this.args.entityId;
  }
}