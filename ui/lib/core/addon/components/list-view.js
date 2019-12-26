import Component from '@ember/component';
import { computed } from '@ember/object';
import { pluralize } from 'ember-inflector';
import layout from '../templates/components/list-view';

/**
 * @module ListView
 * `ListView` components are used in conjuction with `ListItem` for rendering a list.
 *
 * @example
 * ```js
 * <ListView @items={{model}} @itemNoun="role" @paginationRouteName="scope.roles" as |list|>
 *   {{#if list.empty}}
 *     <list.empty @title="No roles here" />
 *   {{else}}
 *     <div>
 *       {{list.item.id}}
 *     </div>
 *   {{/if}}
 * </ListView>
 * ```
 *
 * @param items=null {Array} - An array of items to render as a list
 * @param [itemNoun=null {String}] - A noun to use in the empty state of message and title.
 * @param [message=null {String}] - The message to display within the banner.
 * @yields Object with `item` that is the current item in the loop.
 * @yields If there are no objects in items, then `empty` will be yielded - this is an instance of
 * the EmptyState component.
 * @yields If `item` or `empty` isn't present on the object, the component can still yield a block - this is
 * useful for showing states where there are items but there may be a filter applied that returns an
 * empty set.
 *
 */
export default Component.extend({
  layout,
  tagName: '',
  items: null,
  itemNoun: 'item',
  paginationRouteName: '',
  showPagination: computed('paginationRouteName', 'items.meta{lastPage,total}', function() {
    let meta = this.items.meta;
    return this.paginationRouteName && meta && meta.lastPage > 1 && meta.total > 0;
  }),

  emptyTitle: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `No ${items} yet`;
  }),

  emptyMessage: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `Your ${items} will be listed here. Add your first ${this.get('itemNoun')} to get started.`;
  }),
});
