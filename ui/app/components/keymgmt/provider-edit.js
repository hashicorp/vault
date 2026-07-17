/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { paginate } from 'core/utils/paginate-list';
import { getKeymgmtProviderIcon } from 'vault/utils/keymgmt-provider-utils';
import { SecretsApiKeyManagementListKeysInKmsProviderListEnum } from '@hashicorp/vault-client-typescript';

/**
 * @module KeymgmtProviderEdit
 * ProviderKeyEdit components are used to display KeyMgmt Secrets engine UI for Key items
 *
 * @example
 * ```js
 * <KeymgmtProviderEdit @form={form} @mode="show" />
 * ```
 * @param {object} form - form is the data backing create/edit/show
 * @param {string} mode - mode controls which view is shown on the component - show | create |
 * @param {string} [tab] - Options are "details" or "keys" for the show mode only
 */

export default class KeymgmtProviderEdit extends Component {
  @service api;
  @service capabilities;
  @service router;
  @service flashMessages;

  @tracked modelValidations;
  @tracked invalidFormAlert;
  @tracked isDistributing = false;

  constructor() {
    super(...arguments);
    // key count displayed in details tab and keys are listed in keys tab
    if (this.args.mode === 'show') {
      this.fetchKeys.perform();
    }
  }

  displayFields = ['name', 'provider', 'key_collection', 'keys'];

  label(field) {
    const labels = {
      name: 'Provider name',
      provider: 'Type',
      key_collection: 'Key Vault instance name',
      keys: 'Keys',
    };
    return labels[field] || field;
  }

  providerTypeName(provider) {
    return (
      {
        azurekeyvault: 'Azure Key Vault',
        awskms: 'AWS Key Management Service',
        gcpckms: 'Google Cloud Key Management Service',
      }[provider] || provider
    );
  }

  providerIcon(provider) {
    return getKeymgmtProviderIcon(provider);
  }

  get keyCount() {
    return this.args.form.data.keys?.length || 0;
  }

  get keysValue() {
    if (this.keyCount) {
      return `${this.keyCount} ${this.keyCount > 1 ? 'keys' : 'key'}`;
    }
    return this.args.capabilities?.canListKeys ? 'None' : 'You do not have permission to list keys';
  }

  get breadcrumbs() {
    return [
      {
        label: 'Vault',
        icon: 'vault',
        route: 'vault.cluster.dashboard',
      },
      {
        label: 'Secrets engines',
        route: 'vault.cluster.secrets.backends',
      },
      {
        label: this.args.form.data.backend,
        route: 'vault.cluster.secrets.backend.list-root',
        model: this.args.form.data.backend,
      },
      { label: this.title },
    ];
  }

  get title() {
    if (this.isDistributing) {
      return 'Distribute Key to Provider';
    } else if (this.isShowing) {
      return 'Provider';
    } else if (this.isCreating) {
      return 'Create Provider';
    }

    return 'Update Credentials';
  }

  get subtitle() {
    return this.isShowing ? this.args.form.data.name : '';
  }

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
    const { form } = this.args;
    const { backend, name, provider, key_collection, credentials } = form.data;
    try {
      yield this.api.secrets.keyManagementWriteKmsProvider(name, backend, {
        provider,
        key_collection,
        credentials,
      });

      this.router.transitionTo('vault.cluster.secrets.backend.show', name, {
        queryParams: { itemType: 'provider' },
      });
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }

  @task
  @waitFor
  *fetchKeys(page) {
    const { form, capabilities } = this.args;
    const backend = form.data.backend;
    const providerName = form.data.name;

    if (!capabilities?.canListKeys) {
      form.data.keys = [];
      return;
    }

    try {
      const { keys } = yield this.api.secrets.keyManagementListKeysInKmsProvider(
        providerName,
        backend,
        SecretsApiKeyManagementListKeysInKmsProviderListEnum.TRUE
      );

      const keyNames = keys || [];
      const pathsToFetch = keyNames.map((keyName) =>
        this.capabilities.pathFor('keymgmtKey', { backend, name: keyName })
      );
      const keyCapabilities = yield this.capabilities.fetch(pathsToFetch);

      const keysList = keyNames.map((keyName) => {
        const keyPath = this.capabilities.pathFor('keymgmtKey', { backend, name: keyName });
        return {
          id: keyName,
          name: keyName,
          backend,
          icon: 'key',
          type: 'key',
          canRead: keyCapabilities[keyPath]?.canRead || false,
          canEdit: keyCapabilities[keyPath]?.canUpdate || false,
          canDelete: keyCapabilities[keyPath]?.canDelete || false,
        };
      });

      form.data.keys = paginate(keysList, { page: Number(page) || 1 });
    } catch (error) {
      const { message, status } = yield this.api.parseError(error);
      if (status === 404) {
        form.data.keys = [];
      } else {
        this.flashMessages.danger(message);
      }
    }
  }

  @action
  async onSave(event) {
    event.preventDefault();
    const { isValid, state, invalidFormMessage } = this.args.form.toJSON();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = invalidFormMessage;

    if (isValid) {
      this.modelValidations = null;
      this.saveTask.perform();
    } else {
      this.modelValidations = state;
    }
  }

  @action
  onFieldChange(path) {
    if (path !== 'provider') return;

    // Clear stale validation state on updating Type field so old provider errors do not persist.
    this.modelValidations = null;
    this.invalidFormAlert = null;
  }

  @action
  async onDelete() {
    try {
      const { form, root } = this.args;
      await this.api.secrets.keyManagementDeleteKmsProvider(form.data.name, form.data.backend);
      this.router.transitionTo(root.path, root.model, { queryParams: { tab: 'provider' } });
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }

  @action
  async onDeleteKey(model) {
    try {
      const backend = this.args.form.data.backend;
      const providerName = this.args.form.data.name;
      await this.api.secrets.keyManagementDeleteKeyInKmsProvider(model.id, providerName, backend);
      this.fetchKeys.perform(this.args.form.data.keys?.meta.currentPage || 1);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
