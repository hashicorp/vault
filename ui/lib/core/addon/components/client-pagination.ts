import { action } from '@ember/object';
import { pluralize } from 'ember-inflector';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import config from 'vault/config/environment';
import MutableArray from '@ember/array/mutable';

const { DEFAULT_PAGE_SIZE } = config.APP;

const getSliceIdxForCurrentPage = (pageSize: number, currentPage: number) => {
  const pageIndex = currentPage - 1;
  const first = pageIndex * pageSize;
  const last = first + pageSize;
  return [first, last];
};

interface Args {
  items: MutableArray<unknown>;
  itemNoun?: string;
  filterQuery?: string;
}

/**
 * * @module ClientPagination
 * `ClientPagination` components are used to render a list of records from the store.
 * The list will be automatically paginated when >15 items in the list. If the list is empty
 * an empty state will show with calculated title and message. To add empty state actions,
 * wrap them in `<:emptyActions>` within the component.
 *
 * @example
 * ```js
 * <ClientPagination @itemNoun="role" @filterQuery={{this.filter}} @items={{filter-items this.model this.filter}} >
 *   <:emptyActions>
 *     <LinkTo @route="dashboard">This link renders an an empty state action</LinkTo>
 *   </:emptyActions>
 *   <:item as |role|>
 *     <div data-test-list-item-role={{role.id}}>
 *       I have access to the iterated role info: {{role.name}}
 *     </div>
 *   </:item>
 * </ListView>
 * ```
 *
 */
export default class ClientPaginationComponent extends Component<Args> {
  @tracked currentPage = 1;
  @tracked pageSize = DEFAULT_PAGE_SIZE as number;

  get itemNoun() {
    return this.args.itemNoun || 'item';
  }

  get emptyTitle() {
    const items = pluralize(this.itemNoun);
    return `No ${items} ${this.args.filterQuery ? 'found' : 'yet'}`;
  }

  get emptyMessage() {
    const items = pluralize(this.itemNoun);
    const { filterQuery } = this.args;
    if (filterQuery) {
      return `There are no ${items} matching matching "${filterQuery}"`;
    }
    return `Your ${items} will be listed here. Add your first ${this.itemNoun} to get started.`;
  }

  get shownItems() {
    if (!this.args.items) return [];
    const [first, last] = getSliceIdxForCurrentPage(this.pageSize, this.currentPage);
    return this.args.items.slice(first, last);
  }

  @action handlePageChange(pageNumber: number) {
    this.currentPage = pageNumber;
  }
  @action handlePageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    // when page size changes, go back to first page
    this.currentPage = 1;
  }
}
