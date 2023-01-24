import { action } from '@ember/object';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import Ember from 'ember';

/**
 * Confirm that the user wants to discard unsaved changes before leaving the page.
 * This decorator hooks into the willTransition action. If you override setupController,
 * be sure to set 'model' on the controller to store data or this won't work.
 */
export function withConfirmLeave() {
  return function decorator(SuperClass) {
    if (!Object.prototype.isPrototypeOf.call(Route, SuperClass)) {
      // eslint-disable-next-line
      console.error(
        'withConfirmLeave decorator must be used on instance of ember Route class. Decorator not applied to returned class'
      );
      return SuperClass;
    }
    return class ConfirmLeave extends SuperClass {
      @service store;

      @action
      willTransition(transition) {
        try {
          super.willTransition(...arguments);
        } catch (e) {
          // if the SuperClass doesn't have willTransition
          // defined it will throw an error.
        }
        const model = this.controller.get('model');
        if (model && model.hasDirtyAttributes) {
          if (
            Ember.testing ||
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
