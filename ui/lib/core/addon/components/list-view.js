import Component from '@glimmer/component';
import { pluralize } from 'ember-inflector';

/**
 * @module ListView
 * `ListView` components are used in conjunction with `ListItem` for rendering a list.
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
 * @param {array} [items=null] - An Ember array of items (objects) to render as a list. Because it's an Ember array it has properties like length an meta on it.
 * @param {string} [itemNoun=item] - A noun to use in the empty state of message and title.
 * @param {string} [message=null] - The message to display within the banner.
 * @param {string} [paginationRouteName=''] - The link used in the ListPagination component.
 * @yields {object} Yields the current item in the loop.
 * @yields If there are no objects in items, then `empty` will be yielded - this is an instance of
 * the EmptyState component.
 * @yields If `item` or `empty` isn't present on the object, the component can still yield a block - this is
 * useful for showing states where there are items but there may be a filter applied that returns an
 * empty set.
 *
 */
export default class ListView extends Component {
  get itemNoun() {
    return this.args.itemNoun || 'item';
  }

  get showPagination() {
    let meta = this.args.items.meta;
    return this.args.paginationRouteName && meta && meta.lastPage > 1 && meta.total > 0;
  }

  get emptyTitle() {
    let items = pluralize(this.itemNoun);
    return `No ${items} yet`;
  }

  get emptyMessage() {
    let items = pluralize(this.itemNoun);
    return `Your ${items} will be listed here. Add your first ${this.itemNoun} to get started.`;
  }
}
