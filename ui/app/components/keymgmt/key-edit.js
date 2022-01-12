import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module KeymgmtKeyEdit
 * KeymgmtKeyEdit components are used to display KeyMgmt Secrets engine UI for Key items
 *
 * @example
 * ```js
 * <KeymgmtKeyEdit @model={model} @mode="show" @tab="versions" />
 * ```
 * @param {object} model - model is the data from the store
 * @param {string} [mode=show] - mode controls which view is shown on the component
 * @param {string} [tab=details] - Options are "details" or "versions" for the show mode only
 */

export default class KeymgmtKeyEdit extends Component {
  @service store;
  @service flashMessages;
  @tracked isDeleteModalOpen = false;

  get mode() {
    return this.args.mode || 'show';
  }

  @action
  toggleModal(bool) {
    this.isDeleteModalOpen = bool;
  }

  @action
  createKey(evt) {
    evt.preventDefault();
    this.args.model.save();
  }

  @action
  updateKey(evt) {
    evt.preventDefault();
    this.args.model.save();
  }

  @action
  removeKey(id) {
    // TODO: remove action
    console.log('remove', id);
  }

  @action
  deleteKey(id) {
    // TODO: delete action
    console.log('delete key', id);
    // TODO: Redirect
    this.isDeleteModalOpen = false;
  }

  @action
  rotateKey(id) {
    const backend = this.args.model.get('backend');
    let adapter = this.store.adapterFor('keymgmt/key');
    adapter
      .rotateKey(backend, id)
      .then(() => {
        this.flashMessages.success(`Success: ${id} connection was rotated`);
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors);
      });
  }
}
