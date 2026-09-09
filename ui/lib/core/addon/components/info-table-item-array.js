/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { assert } from '@ember/debug';

/**
 * @module InfoTableItemArray
 * The `InfoTableItemArray` component handles arrays in the info-table-row component.
 * If an array has more than 10 items, then only 5 are displayed and a count of total items is displayed next to the five.
 * If a isLink is true than a link can be set for the use to click on the specific array item
 * Wildcard items are rendered as the raw string passed through the display array.
 *
 * @example
 * <InfoTableItemArray @label="Roles" @displayArray={{array "test-1" "test-2" "test-3"}}  />
 *
 * @param {string} label - used to render lowercased display text for "View all [label]."
 * @param {array} displayArray - The array of data to be displayed. (In InfoTableRow this comes from the `@value` arg.) If the array length > 10, and `@doNotTruncate` is false only 5 will show with a count of the number hidden.
 * @param {boolean} [isLink]  - Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed.
 * @param {string | array} [rootRoute=vault.cluster.secrets.backend.list-root] -  Tells what route the link should go to when selecting "view all". If the route requires more than one dynamic param, insert an array.
 * @param {string | array} [itemRoute=vault.cluster.secrets.backend.show] - Tells what route the link should go to when selecting the individual item. If the route requires more than one dynamic param, insert an array.
 * @param {array} [arrayOptions] - Optional array of preloaded item names used to count wildcard matches.
 * @param {string} [wildcardLabel] - Singular label used when rendering a wildcard match count badge.
 * @param {string} [backend] - To specify which backend to point the link to.
 * @param {boolean} [doNotTruncate=false] - Determines whether to show the View all "roles" link. Otherwise uses the ReadMore component's "See More" toggle
 */
export default class InfoTableItemArray extends Component {
  constructor() {
    super(...arguments);
    assert('@label is required for InfoTableItemArray components', this.args.label);
  }

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

  get wildcardLabel() {
    return this.args.wildcardLabel || '';
  }
}
