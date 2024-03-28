/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module TransformListItem
 * TransformListItem components are used for the list items for the Transform Secret Engines for all but Transformations.
 * This component automatically handles read-only list items if capabilities are not granted or the item is internal only.
 *
 * @example
 * ```js
 * <TransformListItem @item={item} @itemPath="role/my-item" @itemType="role" />
 * ```
 * @param {object} item - item refers to the model item used on the list item partial
 * @param {string} itemPath - usually the id of the item, but can be prefixed with the model type (see transform/role)
 * @param {string} [itemType] - itemType is used to calculate whether an item is readable or
 */

import { computed } from '@ember/object';
import Component from '@ember/component';

export default Component.extend({
  item: null,
  itemPath: '',
  itemType: '',

  isBuiltin: computed('item', 'itemType', function () {
    const item = this.item;
    if (this.itemType === 'alphabet' || this.itemType === 'template') {
      return item.id.startsWith('builtin/');
    }
    return false;
  }),

  backendType: 'transform',
});
