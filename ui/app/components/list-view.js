import Ember from 'ember';
import { pluralize } from 'ember-inflector';
import { task } from 'ember-concurrency';

const { computed, inject } = Ember;
export default Ember.Component.extend({
  flashMessages: inject.service(),
  tagName: '',
  items: null,
  itemNoun: 'item',

  emptyMessage: computed('itemNoun', function() {
    let items = pluralize(this.get('itemNoun'));
    return `There are currently no ${items}`;
  }),

  saveItem: task(function*() {}),
  deleteItem: task(function*(model, successMessage, failureMessage) {
    let flash = this.get('flashMessages');
    try {
      yield model.destroyRecord();
      flash.success(successMessage);
    } catch (e) {
      let errString = e.errors.join(',');
      flash.danger(failureMessage + errString);
      model.rollbackAttributes();
    }
  }),
});
