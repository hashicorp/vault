import Component from '@ember/component';
import { computed } from '@ember/object';
import { pluralize } from 'ember-inflector';
import layout from '../templates/components/list-view';

export default Component.extend({
  layout,
  tagName: '',
  items: null,
  itemNoun: 'item',
  // the dasherized name of a component to render
  // in the EmptyState component if there are no items in items.length
  emptyActions: '',
  showPagination: computed('paginationRouteName', 'items.meta{lastPage,total}', function() {
    return this.paginationRouteName && this.items.meta.lastPage > 1 && this.items.meta.total > 0;
  }),

  paginationRouteName: '',

  emptyTitle: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `No ${items} yet`;
  }),

  emptyMessage: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `Your ${items} will be listed here. Add your first ${this.get('itemNoun')} to get started.`;
  }),
});
