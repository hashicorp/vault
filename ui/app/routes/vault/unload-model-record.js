import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import removeRecord from 'vault/utils/remove-record';
import { get, action } from '@ember/object';

// TODO (unload mixin): change name to UnloadModelRoute when mixin is removed
// removes Ember Data records from the cache when the model changes or you move away from the current route
export default class UnloadModelRecord extends Route {
  @service store;
  modelPath = null; // TODO (unload mixin): make sure migrated mixin routes no longer prepend modelPath with 'model'
  alertUnsavedChanges = false; // defaults to fire deactivate(), if true then calls willTransition()

  get modelToUnload() {
    return !this.modelPath ? this.currentModel : get(this.currentModel, this.modelPath);
  }

  @action
  willTransition(transition) {
    if (this.alertUnsavedChanges && this.modelToUnload?.hasDirtyAttributes) {
      if (
        window.confirm(
          'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
        )
      ) {
        this.modelToUnload.rollbackAttributes();
        return true;
      } else {
        transition.abort();
        return false;
      }
    }
    return true;
  }

  deactivate() {
    super.deactivate(...arguments);
    // error is thrown when you attempt to unload a record that is inFlight (isSaving)
    if (!this.alertUnsavedChanges && this.modelToUnload && !this.modelToUnload.isSaving) {
      removeRecord(this.store, this.modelToUnload);
    }
  }
}
