/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { getOwner } from '@ember/owner';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { singularize } from 'ember-inflector';

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
 * @param {object} urls - Item request paths from the generated-item-urls util.
 */

export default class GeneratedItemList extends Component {
  @service api;
  @service controlGroup;
  @service flashMessages;
  @service router;
  @tracked itemToDelete = null;
  @tracked controlGroupError = null;

  get breadcrumbs() {
    return [
      {
        label: 'Vault',
        route: 'vault.cluster.dashboard',
        icon: 'vault',
      },
      {
        label: 'Auth Methods',
        route: 'vault.cluster.access.methods',
      },
      { label: this.args.methodModel.id },
    ];
  }

  deleteItem = task(
    waitFor(async (item) => {
      const itemType = singularize(this.args.itemType);
      this.controlGroupError = null;
      try {
        await this.api.request.delete(this.args.urls.delete(item.id));
      } catch (error) {
        const { message, response } = await this.api.parseError(error);
        // A delete behind a Control Group needs approval; save the token so a retry after approval unwraps it.
        if (response?.isControlGroupError) {
          this.controlGroup.saveTokenFromError(response);
          this.controlGroupError = this.controlGroup.logFromError(response);
        } else {
          this.flashMessages.danger(`There was an error deleting this ${itemType} ${message}`);
        }
        return;
      }
      this.flashMessages.success(`Successfully deleted ${itemType} ${item.id}.`);
      getOwner(this).lookup(`route:${this.router.currentRouteName}`).refresh();
    })
  );
}
