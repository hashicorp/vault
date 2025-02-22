/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { getOwner } from '@ember/owner';
import { tracked } from '@glimmer/tracking';
import { singularize } from 'ember-inflector';
import { capitalize } from '@ember/string';

/**
 * @module GeneratedItemList
 * The `GeneratedItemList` component lists generated items related to mounts (e.g. groups, roles, users)
 *
 * @example
 * ```js
 * <GeneratedItemList @model={{model}} @itemType={{itemType}} @paths={{this.paths}} @methodModel={{this.methodModel}}/>
 * ```
 *
 * @param {class} model=null - The corresponding item model that is being configured.
 * @param {string} itemType - The type of item displayed.
 * @param {array} paths - Relevant to the link for the LinkTo element.
 * @param {class} methodModel - Model for the particular method selected.
 */

export default class GeneratedItemList extends Component {
  @service router;
  @service pagination;
  @service flashMessages;
  @service api;

  @tracked itemToDelete = null;

  refreshItemList() {
    const route = getOwner(this).lookup(`route:${this.router.currentRouteName}`);
    this.pagination.clearDataset();
    route.refresh();
  }

  @task
  *deleteAuthMethod() {
    const { id, type, listItem, authMethodPath } = this.itemToDelete;
    try {
      const nameKey = type === 'userpass' ? 'username' : 'name';
      const payload = {
        [`${type}MountPath`]: authMethodPath,
        [nameKey]: id,
      };
      const authDeleteMethod = `${type}Delete${capitalize(listItem)}`;

      yield this.api.auth[authDeleteMethod](payload);

      const message = `Successfully deleted ${singularize(listItem)} ${id}.`;
      this.flashMessages.success(message);
      this.refreshItemList();
    } catch (error) {
      const e = (yield error.response?.json()) || error;
      const errString = e.errors?.join(' ') || error.message;
      const message = `There was an error deleting this ${singularize(listItem)}: ${errString}`;
      this.flashMessages.danger(message);
    }
  }
}
