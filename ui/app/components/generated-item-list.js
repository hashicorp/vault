/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import { tracked } from '@glimmer/tracking';

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
  @service store;
  @tracked itemToDelete = null;

  @action
  refreshItemList() {
    const route = getOwner(this).lookup(`route:${this.router.currentRouteName}`);
    this.store.clearAllDatasets();
    route.refresh();
  }
}
