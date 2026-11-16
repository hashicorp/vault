/**
 * Copyright IBM Corp. 2016, 2025
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
  @service api;
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
  async resetConnection(id) {
    const { backend } = this.args.item;
    this.actionRunning = 'reset';
    try {
      await this.api.secrets.databaseResetConnection(id, backend);
      this.flashMessages.success(`Success: ${id} connection was reset`);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    } finally {
      this.actionRunning = null;
    }
  }
  @action
  async rotateRootCred(id) {
    const { backend } = this.args.item;
    this.actionRunning = 'rotateRoot';
    try {
      await this.api.secrets.databaseRotateRootCredentials(id, backend);
      this.flashMessages.success(`Success: ${id} connection was rotated`);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    } finally {
      this.actionRunning = null;
    }
  }
  @action
  async rotateRoleCred(id) {
    const { backend } = this.args.item;
    this.actionRunning = 'rotateRole';
    try {
      await this.api.secrets.databaseRotateStaticRoleCredentials(id, backend);
      this.flashMessages.success(`Success: Credentials for ${id} role were rotated`);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    } finally {
      this.actionRunning = null;
    }
  }
}
