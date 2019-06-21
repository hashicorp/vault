import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import DS from 'ember-data';
import layout from '../templates/components/edit-form';
import { next } from '@ember/runloop';

export default Component.extend({
  layout,
  flashMessages: service(),

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

  // onSave may need values updated in render in a helper - if this
  // is the case, set this value to true
  callOnSaveAfterRender: false,

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
    if (this.callOnSaveAfterRender) {
      next(() => {
        this.get('onSave')({ saveType: method, model });
      });
      return;
    }
    yield this.get('onSave')({ saveType: method, model });
  })
    .drop()
    .withTestWaiter(),

  willDestroy() {
    let model = this.get('model');
    if ((model.get('isDirty') && !model.isDestroyed) || !model.isDestroying) {
      model.rollbackAttributes();
    }
  },
});
