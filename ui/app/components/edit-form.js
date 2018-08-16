import Ember from 'ember';
import { task } from 'ember-concurrency';
import DS from 'ember-data';

const { inject } = Ember;
export default Ember.Component.extend({
  flashMessages: inject.service(),

  // public API
  model: null,
  successMessage: 'Saved!',
  deleteSuccessMessage: 'Deleted!',
  deleteButtonText: 'Delete',
  saveButtonText: 'Save',
  cancelLink: null,

  /*
   * @param Function
   * @public
   *
   * Optional param to call a function upon successfully saving a model
   */
  onSave: () => {},

  save: task(function*(model, options = { method: 'save' }) {
    let { method } = options;
    let messageKey = method === 'save' ? 'successMessage' : 'deleteSuccessMessage';
    try {
      yield model[method]();
    } catch (err) {
      // err will display via model state
      // AdapterErrors are handled by the error-message component
      if (err instanceof DS.AdapterError === false) {
        throw err;
      }
      return;
    }
    this.get('flashMessages').success(this.get(messageKey));
    yield this.get('onSave')({ saveType: method, model });
  }).drop(),

  willDestroy() {
    let model = this.get('model');
    if ((model.get('isDirty') && !model.isDestroyed) || !model.isDestroying) {
      model.rollbackAttributes();
    }
  },
});
