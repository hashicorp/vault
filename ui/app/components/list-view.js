import Component from '@ember/component';
import { computed } from '@ember/object';
import { pluralize } from 'ember-inflector';

export default Component.extend({
  tagName: '',
  items: null,
  itemNoun: 'item',

  emptyTitle: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `No ${items} yet`;
  }),
  emptyMessage: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `Your ${items} will be listed here. Add your first ${this.get('itemNoun')} to get started.`;
  })
});
