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
import { isValidProvider } from 'vault/utils/keymgmt-provider-utils';

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
  @service api;
  @service router;
  @service flashMessages;
  @tracked isDeleteModalOpen = false;
  @tracked isDistributing = false;

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

  displayFields = [
    'name',
    'created',
    'type',
    'deletion_allowed',
    'latest_version',
    'min_enabled_version',
    'last_rotated',
  ];

  label(field) {
    const labels = {
      name: 'Key name',
      created: 'Created',
      type: 'Type',
      deletion_allowed: 'Allow deletion',
      latest_version: 'Current version',
      min_enabled_version: 'Minimum enabled version',
      last_rotated: 'Last rotated',
    };
    return labels[field] || field;
  }

  defaultShown(field) {
    const defaults = {
      min_enabled_version: 'All versions enabled',
      last_rotated: 'Not yet rotated',
    };
    return defaults[field];
  }

  formatDate(field) {
    const dateFields = ['created', 'last_rotated'];
    return dateFields.includes(field) ? 'MMM d yyyy, h:mm:ss aaa' : undefined;
  }

  distributionFields = [
    {
      name: 'name',
      type: 'string',
      label: 'Distributed name',
      subText: 'The name given to the key by the provider.',
    },
    { name: 'purpose', type: 'string', label: 'Key Purpose' },
    { name: 'protection', type: 'string', subText: 'Where cryptographic operations are performed.' },
  ];

  get title() {
    if (this.isDistributing) {
      return 'Distribute key';
    } else if (this.args.mode === 'create') {
      return 'Create key';
    } else if (this.args.mode === 'edit') {
      return 'Edit key';
    }
    return this.args.form.data.name;
  }

  get mode() {
    return this.args.mode || 'show';
  }

  get isMutable() {
    return ['create', 'edit'].includes(this.args.mode);
  }

  get isCreating() {
    return this.args.mode === 'create';
  }

  get hasValidProvider() {
    return isValidProvider(this.args.form?.data?.provider);
  }

  @task
  @waitFor
  *saveKey(evt) {
    evt.preventDefault();
    const { form } = this.args;
    const backend = form.data.backend;
    const name = form.data.name;

    try {
      if (this.isCreating) {
        yield this.api.secrets.keyManagementUpdateKey(name, backend, { type: form.data.type });

        // These fields can only be set after key creation
        try {
          yield this.api.secrets.keyManagementUpdateKey(name, backend, {
            deletion_allowed: form.data.deletion_allowed,
            min_enabled_version: form.data.min_enabled_version || 0,
          });
        } catch (error) {
          this.flashMessages.danger(`Key ${name} was created, but not all settings were saved`);
          this.router.transitionTo(SHOW_ROUTE, name);
          return;
        }
      } else {
        yield this.api.secrets.keyManagementUpdateKey(name, backend, {
          deletion_allowed: form.data.deletion_allowed,
          min_enabled_version: form.data.min_enabled_version,
        });
      }

      this.router.transitionTo(SHOW_ROUTE, name);
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }

  @task
  @waitFor
  *removeKey() {
    const { form } = this.args;
    const backend = form.data.backend;
    const name = form.data.name;
    const provider = form.data.provider;

    if (!this.hasValidProvider) {
      this.flashMessages.danger('Cannot remove key: invalid provider');
      return;
    }

    try {
      yield this.api.secrets.keyManagementDeleteKeyInKmsProvider(name, provider, backend);
      this.flashMessages.success('Key has been successfully removed from provider');
      this.router.refresh();
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }

  @action
  deleteKey() {
    const { form } = this.args;
    const backend = form.data.backend;
    const name = form.data.name;

    this.api.secrets
      .keyManagementDeleteKey(name, backend)
      .then(() => {
        this.router.transitionTo(LIST_ROOT_ROUTE, backend);
      })
      .catch(async (e) => {
        const { message } = await this.api.parseError(e);
        this.flashMessages.danger(message);
      });
  }

  @task
  @waitFor
  *rotateKey() {
    const { form } = this.args;
    const name = form.data.name;
    const backend = form.data.backend;

    try {
      yield this.api.secrets.keyManagementRotateKey(name, backend);
      this.flashMessages.success(`Success: ${name} key was rotated`);
      this.router.refresh();
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
