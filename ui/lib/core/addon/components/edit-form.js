/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import layout from '../templates/components/edit-form';
import { next } from '@ember/runloop';
import { waitFor } from '@ember/test-waiters';

export default Component.extend({
  layout,
  flashMessages: service(),

  // internal validations
  invalidFormAlert: '',
  modelValidations: null,
  // public API
  model: null,
  successMessage: 'Saved!',
  deleteSuccessMessage: 'Deleted!',
  deleteButtonText: 'Delete',
  saveButtonText: 'Save',
  cancelButtonText: 'Cancel',
  cancelLink: null,
  flashEnabled: true,
  includeBox: true,

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

  checkModelValidity(model) {
    if (model.validate) {
      const { isValid, state, invalidFormMessage } = model.validate();
      this.set('modelValidations', state);
      this.set('invalidFormAlert', invalidFormMessage);
      return isValid;
    }
    // no validations on model; return true
    return true;
  },

  save: task(
    waitFor(function* (model, options = { method: 'save' }) {
      const { method } = options;
      const messageKey = method === 'save' ? 'successMessage' : 'deleteSuccessMessage';
      if (method === 'save' && !this.checkModelValidity(model)) {
        // if saving and model invalid, don't continue
        return;
      }
      try {
        yield model[method]();
      } catch (err) {
        // err will display via model state
        // AdapterErrors are handled by the error-message component
        if (err instanceof AdapterError === false) {
          throw err;
        }
        return;
      }
      if (this.flashEnabled) {
        this.flashMessages.success(this[messageKey]);
      }
      if (this.callOnSaveAfterRender) {
        next(() => {
          this.onSave({ saveType: method, model });
        });
        return;
      }
      this.onSave({ saveType: method, model });
    })
  ).drop(),

  willDestroy() {
    // components are torn down after store is unloaded and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    const { model } = this;
    if (noTeardown && model && model.isDirty && !model.isDestroyed && !model.isDestroying) {
      model.rollbackAttributes();
    }
    this._super(...arguments);
  },
});
