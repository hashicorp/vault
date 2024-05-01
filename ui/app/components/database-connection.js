/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { waitFor } from '@ember/test-waiters';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

const getErrorMessage = (errors) => {
  let errorMessage = errors?.join('. ') || 'Something went wrong. Check the Vault logs for more information.';
  if (errorMessage.indexOf('failed to verify') >= 0) {
    errorMessage =
      'There was a verification error for this connection. Check the Vault logs for more information.';
  }
  return errorMessage;
};

export default class DatabaseConnectionEdit extends Component {
  @service store;
  @service router;
  @service flashMessages;

  @tracked
  showPasswordField = false; // used for edit mode

  @tracked
  showSaveModal = false; // used for create mode

  rotateCredentials(backend, name) {
    const adapter = this.store.adapterFor('database/connection');
    return adapter.rotateRootCredentials(backend, name);
  }

  transitionToRoute() {
    return this.router.transitionTo(...arguments);
  }

  @action
  updateShowPassword(showForm) {
    this.showPasswordField = showForm;
    if (!showForm) {
      // unset password if hidden
      this.args.model.password = undefined;
    }
  }

  @action
  updatePassword(attr, evt) {
    const value = evt.target.value;
    this.args.model[attr] = value;
  }

  @action
  async handleCreateConnection(evt) {
    evt.preventDefault();
    const secret = this.args.model;
    secret
      .save()
      .then(() => {
        this.showSaveModal = true;
      })
      .catch((e) => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }

  @action
  continueWithoutRotate() {
    this.showSaveModal = false;
    const { name } = this.args.model;
    this.transitionToRoute(SHOW_ROUTE, name);
  }

  @action
  @waitFor
  async continueWithRotate() {
    this.showSaveModal = false;
    const { backend, name } = this.args.model;
    try {
      await this.rotateCredentials(backend, name);
      this.flashMessages.success(`Successfully rotated root credentials for connection "${name}"`);
      this.transitionToRoute(SHOW_ROUTE, name);
    } catch (e) {
      this.flashMessages.danger(`Error rotating root credentials: ${e.errors}`);
      this.transitionToRoute(SHOW_ROUTE, name);
    }
  }

  @action
  handleUpdateConnection(evt) {
    evt.preventDefault();
    const secret = this.args.model;
    const secretId = secret.name;
    secret
      .save()
      .then(() => {
        this.transitionToRoute(SHOW_ROUTE, secretId);
      })
      .catch((e) => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }

  @action
  delete(evt) {
    evt.preventDefault();
    const secret = this.args.model;
    const backend = secret.backend;
    secret.destroyRecord().then(() => {
      this.transitionToRoute(LIST_ROOT_ROUTE, backend);
    });
  }

  @action
  reset() {
    const { name, backend } = this.args.model;
    const adapter = this.store.adapterFor('database/connection');
    adapter
      .resetConnection(backend, name)
      .then(() => {
        // TODO: Why isn't the confirmAction closing?
        this.flashMessages.success('Successfully reset connection');
      })
      .catch((e) => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }

  @action
  rotate() {
    const { name, backend } = this.args.model;
    this.rotateCredentials(backend, name)
      .then(() => {
        // TODO: Why isn't the confirmAction closing?
        this.flashMessages.success('Successfully rotated credentials');
      })
      .catch((e) => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }
}
