/**
 * @module TransformListItem
 * TransformListItem components are used to...
 *
 * @example
 * ```js
 * <TransformListItem @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
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

  itemViewable: computed('item', 'itemType', function() {
    const item = this.get('item');
    if (this.itemType === 'alphabet' || this.itemType === 'template') {
      return !item.get('id').startsWith('builtin/');
    }
    return true;
  }),

  backendType: 'transform',
});
