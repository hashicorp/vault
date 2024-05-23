/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';
import { tracked } from '@glimmer/tracking';
import keys from 'core/utils/key-codes';

/**
 * @module Roles
 * RolesPage component is a child component to show list of roles.
 * It also handles the filtering actions of roles.
 *
 * @param {array} roles - array of roles
 * @param {boolean} promptConfig - whether or not to display config cta
 * @param {string} filterValue - value of queryParam pageFilter
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */
export default class RolesPageComponent extends Component {
  @service flashMessages;
  @service router;
  @tracked query;
  @tracked roleToDelete = null;

  constructor() {
    super(...arguments);
    this.query = this.args.filterValue;
  }

  get mountPoint() {
    return getOwner(this).mountPoint;
  }

  navigate(pageFilter) {
    const route = `${this.mountPoint}.roles.index`;
    const args = [route, { queryParams: { pageFilter: pageFilter || null } }];
    this.router.transitionTo(...args);
  }

  @action
  handleKeyDown(event) {
    if (event.keyCode === keys.ESC) {
      // On escape, transition to roles index route.
      this.navigate();
    }
    // ignore all other key events
  }

  @action handleInput(evt) {
    this.query = evt.target.value;
  }

  @action
  handleSearch(evt) {
    evt.preventDefault();
    this.navigate(this.query);
  }

  @action
  async onDelete(model) {
    try {
      const message = `Successfully deleted role ${model.name}`;
      await model.destroyRecord();
      this.flashMessages.success(message);
    } catch (error) {
      const message = errorMessage(error, 'Error deleting role. Please try again or contact support');
      this.flashMessages.danger(message);
    } finally {
      this.roleToDelete = null;
    }
  }
}
