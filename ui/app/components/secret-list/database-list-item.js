/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module DatabaseListItem
 * DatabaseListItem components are used for the list items for the Database Secret Engines for Roles.
 * This component automatically handles read-only list items if capabilities are not granted or the item is internal only.
 *
 * @example
 * ```js
 * <DatabaseListItem @item={item} />
 * ```
 * @param {object} item - item refers to the model item used on the list item partial
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class DatabaseListItem extends Component {
  @tracked roleType = '';
  @tracked actionRunning = null;
  @service store;
  @service flashMessages;

  get keyTypeValue() {
    const item = this.args.item;
    // basing this on path in case we want to remove 'type' later
    if (item.path === 'roles') {
      return 'dynamic';
    } else if (item.path === 'static-roles') {
      return 'static';
    } else {
      return '';
    }
  }

  @action
  resetConnection(id) {
    const { backend } = this.args.item;
    const adapter = this.store.adapterFor('database/connection');
    this.actionRunning = 'reset';
    adapter
      .resetConnection(backend, id)
      .then(() => {
        this.flashMessages.success(`Success: ${id} connection was reset`);
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors);
      })
      .finally(() => (this.actionRunning = null));
  }
  @action
  rotateRootCred(id) {
    const { backend } = this.args.item;
    const adapter = this.store.adapterFor('database/connection');
    this.actionRunning = 'rotateRoot';
    adapter
      .rotateRootCredentials(backend, id)
      .then(() => {
        this.flashMessages.success(`Success: ${id} connection was rotated`);
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors);
      })
      .finally(() => (this.actionRunning = null));
  }
  @action
  rotateRoleCred(id) {
    const { backend } = this.args.item;
    const adapter = this.store.adapterFor('database/credential');
    this.actionRunning = 'rotateRole';
    adapter
      .rotateRoleCredentials(backend, id)
      .then(() => {
        this.flashMessages.success(`Success: Credentials for ${id} role were rotated`);
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors);
      })
      .finally(() => (this.actionRunning = null));
  }
}
