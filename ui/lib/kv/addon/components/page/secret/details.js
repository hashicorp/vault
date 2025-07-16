/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { next } from '@ember/runloop';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { isDeleted } from 'kv/utils/kv-deleted';
import { isAdvancedSecret } from 'core/utils/advanced-secret';

/**
 * @module KvSecretDetails renders the key/value data of a KV secret.
 * It also renders a dropdown to display different versions of the secret.
 * <Page::Secret::Details
 * @backend={{this.model.backend}}
 * @breadcrumbs={{this.breadcrumbs}}
 * @canReadData={{this.model.canReadData}}
 * @canReadMetadata={{this.model.canReadMetadata}}
 * @canUpdateData={{this.model.canUpdateData}}
 * @isPatchAllowed={{this.model.isPatchAllowed}}
 * @metadata={{this.model.metadata}}
 * @path={{this.model.path}}
 * @secret={{this.model.secret}}
 * />
 *
 * @param {string} backend - path where kv engine is mounted
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * @param {boolean} canReadData - if true and the secret is not destroyed/deleted the copy secret dropdown renders
 * @param {boolean} canReadMetadata - if true it renders the kv select version dropdown in the toolbar and "Version History" tab
 * @param {boolean} canUpdateData - if true it renders "Create new version" toolbar action
 * @param {boolean} isPatchAllowed - if true it renders "Patch latest version" toolbar action. True when: (1) the version is enterprise, (2) a user has "patch" secret + "read" subkeys capabilities, (3) latest secret version is not deleted or destroyed
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {string} path - path of kv secret 'my/secret' used as the title for the KV page header
 * @param {model} secret - Ember data model: 'kv/data'
 */

export default class KvSecretDetails extends Component {
  @service flashMessages;
  @service('app-router') router;
  @service store;

  @tracked showJsonView = false;
  @tracked wrappedData = null;
  @tracked syncStatus = null; // array of association sync status info by destination

  constructor() {
    super(...arguments);
    this.fetchSyncStatus.perform();
    this.originalSecret = JSON.stringify(this.args.secret.secretData || {});
    if (isAdvancedSecret(this.originalSecret)) {
      // Default to JSON view if advanced
      this.showJsonView = true;
    }
  }

  @action
  closeVersionMenu(dropdown) {
    // strange issue where closing dropdown triggers full transition (which redirects to auth screen in production)
    // closing dropdown in next tick of run loop fixes it
    next(() => dropdown.actions.close());
  }

  @action
  clearWrappedData() {
    this.wrappedData = null;
  }

  @task
  @waitFor
  *wrapSecret() {
    const { backend, path } = this.args.secret;
    const adapter = this.store.adapterFor('kv/data');
    try {
      const { token } = yield adapter.fetchWrapInfo({ backend, path, wrapTTL: 1800 });
      if (!token) throw 'No token';
      this.wrappedData = token;
      this.flashMessages.success('Secret successfully wrapped!');
    } catch (error) {
      this.flashMessages.danger('Could not wrap secret.');
    }
  }

  @task
  @waitFor
  *fetchSyncStatus() {
    const { backend, path } = this.args.secret;
    const syncAdapter = this.store.adapterFor('sync/association');
    try {
      this.syncStatus = yield syncAdapter.fetchSyncStatus({ mount: backend, secretName: path });
    } catch (e) {
      // silently error
    }
  }

  @action
  async undelete() {
    const { secret } = this.args;
    try {
      await secret.destroyRecord({
        adapterOptions: { deleteType: 'undelete', deleteVersions: this.version },
      });
      this.flashMessages.success(`Successfully undeleted ${secret.path}.`);
      this.transition();
    } catch (err) {
      this.flashMessages.danger(
        `There was a problem undeleting ${secret.path}. Error: ${err.errors?.join(' ')}.`
      );
    }
  }

  @action
  async handleDestruction(type) {
    const { secret } = this.args;
    try {
      await secret.destroyRecord({ adapterOptions: { deleteType: type, deleteVersions: this.version } });
      const verb = type.includes('delete') ? 'deleted' : 'destroyed';
      this.flashMessages.success(`Successfully ${verb} Version ${this.version} of ${secret.path}.`);
      this.transition();
    } catch (err) {
      const verb = type.includes('delete') ? 'deleting' : 'destroying';
      this.flashMessages.danger(
        `There was a problem ${verb} Version ${this.version} of ${secret.path}. Error: ${err.errors.join(
          ' '
        )}.`
      );
    }
  }

  transition() {
    // transition to the overview to prevent automatically reading sensitive secret data
    this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.index');
  }

  get version() {
    return (
      this.args.secret?.version ||
      this.router.currentRoute.queryParams?.version ||
      this.args.metadata?.sortedVersions[0].version
    );
  }

  get hideHeaders() {
    return this.showJsonView || this.emptyState;
  }

  get versionState() {
    const { secret, metadata } = this.args;
    if (secret.failReadErrorCode !== 403) {
      return secret.state;
    }
    // If the user can't read secret data, get the current version
    // state from metadata versions
    if (metadata?.sortedVersions) {
      const version = this.version;
      const meta = version
        ? metadata.sortedVersions.find((v) => v.version == version)
        : metadata.sortedVersions[0];
      if (meta?.destroyed) {
        return 'destroyed';
      }
      if (isDeleted(meta?.deletion_time)) {
        return 'deleted';
      }
      if (meta?.created_time) {
        return 'created';
      }
    }
    return '';
  }

  get showUndelete() {
    const { secret } = this.args;
    if (secret.canUndelete) {
      return this.versionState === 'deleted';
    }
    return false;
  }

  get showDelete() {
    const { secret } = this.args;
    if (secret.canDeleteVersion || secret.canDeleteLatestVersion) {
      return this.versionState === 'created' || this.versionState === '';
    }
    return false;
  }

  get showDestroy() {
    const { secret } = this.args;
    if (secret.canDestroyVersion) {
      return this.versionState !== 'destroyed' && this.version;
    }
    return false;
  }

  get emptyState() {
    if (!this.args.canReadData) {
      return {
        title: 'You do not have permission to read this secret',
        message:
          'Your policies may permit you to write a new version of this secret, but do not allow you to read its current contents.',
      };
    }
    // only destructure if we can read secret data
    const { version, destroyed, isSecretDeleted } = this.args.secret;
    if (destroyed) {
      return {
        title: `Version ${version} of this secret has been permanently destroyed`,
        message: `A version that has been permanently deleted cannot be restored. ${
          this.args.canReadMetadata
            ? ' You can view other versions of this secret in the Version History tab above.'
            : ''
        }`,
        link: '/vault/docs/secrets/kv/kv-v2',
      };
    }
    if (isSecretDeleted) {
      return {
        title: `Version ${version} of this secret has been deleted`,
        message: `This version has been deleted but can be undeleted. ${
          this.args.canReadMetadata
            ? 'View other versions of this secret by clicking the Version History tab above.'
            : ''
        }`,
        link: '/vault/docs/secrets/kv/kv-v2',
      };
    }
    return false;
  }
}
