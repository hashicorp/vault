import Route from '@ember/routing/route';
import { action } from '@ember/object';
import removeRecord from 'vault/utils/remove-record';
import { inject as service } from '@ember/service';

export default class UnloadModelRecord extends Route {
  @service store;

  unloadModel() {
    const model = this.currentModel;
    // error is thrown when you attempt to unload a record that is inFlight (isSaving)
    if (!model || !model.unloadRecord || model.isSaving) {
      return;
    }
    removeRecord(this.store, model);
    model.destroy();
  }

  resetController(controller) {
    // it's important to unset the model on the controller since controllers are singletons
    controller.model = null;
  }

  @action
  willTransition(transition) {
    if (transition.to.name.includes('edit')) {
      // don't unload model if about to edit it
      return true;
    }
    if (this.currentModel.hasDirtyAttributes) {
      if (
        window.confirm(
          'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
        )
      ) {
        this.unloadModel();
        return true;
      } else {
        transition.abort();
        return false;
      }
    }
    this.unloadModel();
    return true;
  }
}
