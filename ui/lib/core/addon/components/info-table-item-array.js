/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';

/**
 * @module InfoTableItemArray
 * The `InfoTableItemArray` component handles arrays in the info-table-row component.
 * If an array has more than 10 items, then only 5 are displayed and a count of total items is displayed next to the five.
 * If a isLink is true than a link can be set for the use to click on the specific array item
 * If a wildcard is a potential variable in the string of an item, then you can use the modelType and wildcardLabel parameters to
 * return a wildcard count similar to what is done in the searchSelect component.
 *
 * @example
 * ```js
 * <InfoTableItemArray
 * @label="Roles"
 * @displayArray={{['test-1','test-2','test-3']}}
 * @isLink={{true}}
 * @rootRoute="vault.cluster.secrets.backend.list-root"
 * @itemRoute="vault.cluster.secrets.backend.show"
 * @modelType="transform/role"
 * @queryParam="role"
 * @backend="transform"
 * ```
 *
 * @param {string} label - used to render lowercased display text for "View all <label>."
 * @param {array} displayArray - The array of data to be displayed. (In InfoTableRow this comes from the @value arg.) If the array length > 10, and @doNotTruncate is false only 5 will show with a count of the number hidden.
 * @param {boolean} [isLink]  - Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed.
 * @param {string || array} [rootRoute="vault.cluster.secrets.backend.list-root"] -  Tells what route the link should go to when selecting "view all". If the route requires more than one dynamic param, insert an array.
 * @param {string || array} [itemRoute=vault.cluster.secrets.backend.show] - Tells what route the link should go to when selecting the individual item. If the route requires more than one dynamic param, insert an array.
 * @param {string} [modelType]  - Tells which model you want to query and set allOptions.  Used in conjunction with the the isLink.
 * @param {string} [wildcardLabel]  - when you want the component to return a count on the model for options returned when using a wildcard you must provide a label of the count e.g. role.  Should be singular.
 * @param {string} [backend] - To specify which backend to point the link to.
 * @param {boolean} [doNotTruncate=false] - Determines whether to show the View all "roles" link. Otherwise uses the ReadMore component's "See More" toggle
 * @param {boolean} [renderItemName=false] - If true renders the item name instead of its id
 */
export default class InfoTableItemArray extends Component {
  @service store;
  @tracked allOptions = null;
  @tracked itemNameById; // object is only created if renderItemName=true
  @tracked fetchComplete = false;

  get rootRoute() {
    return this.args.rootRoute || 'vault.cluster.secrets.backend.list-root';
  }

  get itemRoute() {
    return this.args.itemRoute || 'vault.cluster.secrets.backend.show';
  }

  get doNotTruncate() {
    return this.args.doNotTruncate || false;
  }

  get displayArrayTruncated() {
    const { displayArray } = this.args;
    if (!displayArray) return null;
    if (displayArray.length >= 10 && !this.args.doNotTruncate) {
      // if array greater than 10 in length only display the first 5
      return displayArray.slice(0, 5);
    }
    return displayArray;
  }

  @action async fetchOptions() {
    if (this.args.isLink && this.args.modelType) {
      const queryOptions = this.args.backend ? { backend: this.args.backend } : {};

      const modelRecords = await this.store.query(this.args.modelType, queryOptions).catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          return null;
        }
      });

      this.allOptions = modelRecords ? modelRecords.map((record) => record.id) : null;
      if (this.args.renderItemName && modelRecords) {
        modelRecords.forEach(({ id, name }) => {
          // create key/value pair { item-id: item-name } for each record
          this.itemNameById = { ...this.itemNameById, [id]: name };
        });
      }
    }
    this.fetchComplete = true;
  }
}
