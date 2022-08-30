import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module KeymgmtProviderEdit
 * ProviderKeyEdit components are used to display KeyMgmt Secrets engine UI for Key items
 *
 * @example
 * ```js
 * <KeymgmtProviderEdit @model={model} @mode="show" />
 * ```
 * @param {object} model - model is the data from the store
 * @param {string} mode - mode controls which view is shown on the component - show | create |
 * @param {string} [tab] - Options are "details" or "keys" for the show mode only
 */

export default class KeymgmtProviderEdit extends Component {
  @service router;
  @service flashMessages;

  constructor() {
    super(...arguments);
    // key count displayed in details tab and keys are listed in keys tab
    if (this.args.mode === 'show') {
      this.fetchKeys.perform();
    }
  }

  @tracked modelValidations;

  get isShowing() {
    return this.args.mode === 'show';
  }
  get isCreating() {
    return this.args.mode === 'create';
  }
  get viewingKeys() {
    return this.args.tab === 'keys';
  }

  @task
  @waitFor
  *saveTask() {
    const { model } = this.args;
    try {
      yield model.save();
      this.router.transitionTo('vault.cluster.secrets.backend.show', model.id, {
        queryParams: { itemType: 'provider' },
      });
    } catch (error) {
      this.flashMessages.danger(error.errors.join('. '));
    }
  }
  @task
  @waitFor
  *fetchKeys(page = 1) {
    try {
      yield this.args.model.fetchKeys(page);
    } catch (error) {
      this.flashMessages.danger(error.errors.join('. '));
    }
  }

  @action
  async onSave(event) {
    event.preventDefault();
    const { isValid, state } = await this.args.model.validate();
    if (isValid) {
      this.modelValidations = null;
      this.saveTask.perform();
    } else {
      this.modelValidations = state;
    }
  }
  @action
  async onDelete() {
    try {
      const { model, root } = this.args;
      await model.destroyRecord();
      this.router.transitionTo(root.path, root.model, { queryParams: { tab: 'provider' } });
    } catch (error) {
      this.flashMessages.danger(error.errors.join('. '));
    }
  }
  @action
  async onDeleteKey(model) {
    try {
      await model.destroyRecord();
      this.args.model.keys.removeObject(model);
    } catch (error) {
      this.flashMessages.danger(error.errors.join('. '));
    }
  }
}
