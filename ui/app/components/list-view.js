import Ember from 'ember';
import { pluralize } from 'ember-inflector';
import { task } from 'ember-concurrency';

const { computed } = Ember;
export default Ember.Component.extend({
  tagName: '',
  items: null,
  itemNoun: 'item',

  emptyMessage: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `There are currently no ${items}`;
  }),
});
