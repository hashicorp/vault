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

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';
export default class KeymgmtKeyEdit extends Component {
  @service store;
  @service router;
  @service flashMessages;
  @tracked isDeleteModalOpen = false;

  get mode() {
    return this.args.mode || 'show';
  }

  get keyAdapter() {
    return this.store.adapterFor('keymgmt/key');
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
    const name = this.args.model.name;
    this.args.model
      .save()
      .then(() => {
        this.router.transitionTo(SHOW_ROUTE, name);
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors.join('. '));
      });
  }

  @action
  removeKey(id) {
    // TODO: remove action
    console.log('remove', id);
  }

  @action
  deleteKey() {
    const secret = this.args.model;
    const backend = secret.backend;
    console.log({ secret });
    secret
      .destroyRecord()
      .then(() => {
        try {
          this.router.transitionTo(LIST_ROOT_ROUTE, backend, { queryParams: { tab: 'key' } });
        } catch (e) {
          console.debug(e);
        }
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors?.join('. '));
      });
  }

  @action
  rotateKey(id) {
    const backend = this.args.model.get('backend');
    const adapter = this.keyAdapter;
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
