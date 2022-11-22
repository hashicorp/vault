import { action } from '@ember/object';
import { inject as service } from '@ember/service';

/**
 * Confirm that the user wants to discard unsaved changes before leaving the page.
 * This decorator hooks into the willTransition action. If you override setupController,
 * be sure to set 'model' on the controller to store data or this won't work.
 */
export function withConfirmLeave() {
  return function decorator(SuperClass) {
    return class ModelValidations extends SuperClass {
      @service store;

      @action
      willTransition(transition) {
        const model = this.controller.get('model');
        if (model && model.hasDirtyAttributes) {
          if (
            window.confirm(
              'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
            )
          ) {
            // error is thrown when you attempt to unload a record that is inFlight (isSaving)
            if (!model || !model.unloadRecord || model.isSaving) {
              return;
            }
            model.rollbackAttributes();
            model.destroy();
            return true;
          } else {
            transition.abort();
            return false;
          }
        }
      }
    };
  };
}
