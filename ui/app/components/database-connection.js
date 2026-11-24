/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { waitFor } from '@ember/test-waiters';
import { task } from 'ember-concurrency';
import { getShowAttrs, isAvailablePlugin } from 'vault/forms/database/connection';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

const getErrorMessage = (message) => {
  if (!message) {
    return 'Something went wrong. Check the Vault logs for more information.';
  }
  if (message.indexOf('failed to verify') >= 0) {
    return 'There was a verification error for this connection. Check the Vault logs for more information.';
  }
  return message;
};

export default class DatabaseConnectionEdit extends Component {
  @service api;
  @service router;
  @service flashMessages;
  @service version;

  @tracked
  showPasswordField = false; // used for edit mode

  @tracked
  showSaveModal = false; // used for create mode

  @tracked
  error = '';

  breadcrumbs = [
    { label: 'Vault', text: 'Vault', icon: 'vault', path: 'vault.cluster.dashboard' },
    { text: 'Secrets engines', path: 'vault.cluster.secrets.backends' },
    this.args.root,
    { label: this.title, text: this.title },
  ];

  get title() {
    if (this.args?.mode === 'create') {
      return 'Create connection';
    } else if (this.args?.mode === 'edit') {
      return 'Edit connection';
    } else {
      return this.args?.model?.name;
    }
  }

  get isAvailablePlugin() {
    return isAvailablePlugin(this.args.model.plugin_name);
  }

  // fields displayed on the show view, derived from the selected plugin
  get showAttrs() {
    return getShowAttrs(this.args.model.plugin_name, this.version.isEnterprise);
  }

  // Saves the connection from the form. Plugin-specific params are sent in the payload.
  async writeConnection() {
    const { name, backend } = this.args.model;
    const { data } = this.args.model.toJSON();
    // Uses api.request instead of api.secrets.databaseConfigureConnection because the generated request
    // serializer (DatabaseConfigureConnectionRequestToJSON) only keeps the generic connection params and would
    // drop plugin-specific ones such as connection_url, username and password.
    await this.api.request.post(`/${backend}/config/${name}`, data);
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

  handleCreateConnection = task(
    waitFor(async (evt) => {
      evt.preventDefault();
      try {
        await this.writeConnection();
        this.showSaveModal = true;
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.flashMessages.danger(getErrorMessage(message));
      }
    })
  );

  @action
  continueWithoutRotate() {
    if (this.continueWithRotate.isRunning) return;
    this.showSaveModal = false;
    const { name } = this.args.model;
    this.transitionToRoute(SHOW_ROUTE, name);
  }

  continueWithRotate = task(
    waitFor(async () => {
      const { backend, name } = this.args.model;
      try {
        await this.api.secrets.databaseRotateRootCredentials(name, backend);
        this.flashMessages.success(`Successfully rotated root credentials for connection "${name}"`);
        this.transitionToRoute(SHOW_ROUTE, name);
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.flashMessages.danger(`Error rotating root credentials: ${message}`);
        this.transitionToRoute(SHOW_ROUTE, name);
      } finally {
        this.showSaveModal = false;
      }
    })
  );

  @action
  async handleUpdateConnection(evt) {
    evt.preventDefault();
    const { name } = this.args.model;
    try {
      await this.writeConnection();
      this.transitionToRoute(SHOW_ROUTE, name);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(getErrorMessage(message));
    }
  }

  @action
  async delete(evt) {
    evt.preventDefault();
    const { name, backend } = this.args.model;
    try {
      await this.api.secrets.databaseDeleteConnectionConfiguration(name, backend);
      this.transitionToRoute(LIST_ROOT_ROUTE, backend);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.error = message;
    }
  }

  @action
  async reset() {
    const { name, backend } = this.args.model;
    try {
      await this.api.secrets.databaseResetConnection(name, backend);
      this.flashMessages.success('Successfully reset connection');
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(getErrorMessage(message));
    }
  }

  @action
  async rotate() {
    const { name, backend } = this.args.model;
    try {
      await this.api.secrets.databaseRotateRootCredentials(name, backend);
      this.flashMessages.success('Successfully rotated credentials');
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(getErrorMessage(message));
    }
  }
}
