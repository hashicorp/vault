import Ember from 'ember';
import { task } from 'ember-concurrency';
import DS from 'ember-data';

const { computed, inject } = Ember;
export default Ember.Component.extend({
  flashMessages: inject.service(),
  model: null,
  successMessage: 'Saved!',
  deleteSuccessMessage: 'Deleted!',
  deleteButtonText: 'Delete',

  /*
   * @param Function
   * @public
   *
   * Optional param to call a function upon successfully saving an entity
   */
  onSave: () => {},

  save: task(function*() {
    let model = this.get('model');

    try {
      yield model.save();
    } catch (err) {
      // err will display via model state
      // AdapterErrors are handled by the error-message component
      if (err instanceof DS.AdapterError === false) {
        throw err;
      }
      return;
    }
    this.get('flashMessages').success(this.get('successMessage'));
    yield this.get('onSave')({ saveType: 'save', model });
  }).drop(),

  willDestroy() {
    let model = this.get('model');
    if ((model.get('isDirty') && !model.isDestroyed) || !model.isDestroying) {
      model.rollbackAttributes();
    }
  },

  actions: {
    deleteItem(model) {
      let flash = this.get('flashMessages');
      model.destroyRecord().then(() => {
        flash.success(this.get('deleteSuccessMessage'));
        return this.get('onSave')({ saveType: 'delete', model });
      });
    },
  },
});
