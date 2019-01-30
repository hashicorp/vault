import Component from '@ember/component';
import { computed } from '@ember/object';
import { pluralize } from 'ember-inflector';

export default Component.extend({
  tagName: '',
  items: null,
  itemNoun: 'item',
  // the dasherized name of a component to render
  // in the EmptyState component if there are no items in items.length
  emptyActions: '',

  emptyTitle: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `No ${items} yet`;
  }),

  emptyMessage: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `Your ${items} will be listed here. Add your first ${this.get('itemNoun')} to get started.`;
  }),
});
