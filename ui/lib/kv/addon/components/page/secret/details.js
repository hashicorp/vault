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
import sortedVersions from 'kv/helpers/sorted-versions';
import isDeleted from 'kv/helpers/is-deleted';
import { isAdvancedSecret } from 'core/utils/advanced-secret';

/**
 * @module KvSecretDetails renders the key/value data of a KV secret.
 * It also renders a dropdown to display different versions of the secret.
 * <Page::Secret::Details
 *   @backend={{this.model.backend}}
 *   @breadcrumbs={{this.breadcrumbs}}
 *   @capabilities={{this.model.capabilities}}
 *   @isPatchAllowed={{this.model.isPatchAllowed}}
 *   @metadata={{this.model.metadata}}
 *   @path={{this.model.path}}
 *   @secret={{this.model.secret}}
 * />
 *
 * @param {string} backend - path where kv engine is mounted
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * @param {object} capabilities - capabilities for data, metadata, subkeys, delete and undelete paths
 * @param {boolean} isPatchAllowed - if true it renders "Patch latest version" toolbar action. True when: (1) the version is enterprise, (2) a user has "patch" secret + "read" subkeys capabilities, (3) latest secret version is not deleted or destroyed
 * @param {object} metadata - response object from /secret/metadata/path endpoint
 * @param {string} path - path of kv secret 'my/secret' used as the title for the KV page header
 * @param {object} secret - data and metadata objects from kvV2Read response - { secretData: data, ...metadata }
 */

export default class KvSecretDetails extends Component {
  @service flashMessages;
  @service('app-router') router;
  @service api;

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
    try {
      const { secretData: data, ...metadata } = this.args.secret;
      const { wrap_info } = yield this.api.sys.wrap(
        { data, metadata },
        this.api.buildHeaders({ wrap: 1800 })
      );
      if (!wrap_info.token) throw 'No token';
      this.wrappedData = wrap_info.token;
      this.flashMessages.success('Secret successfully wrapped!');
    } catch (error) {
      this.flashMessages.danger('Could not wrap secret.');
    }
  }

  @task
  @waitFor
  *fetchSyncStatus() {
    try {
      const { backend: mount, path: secret_name } = this.args;
      const { associated_destinations } = yield this.api.sys.systemReadSyncAssociationsDestinations(
        (context) => this.api.addQueryParams(context, { mount, secret_name })
      );
      this.syncStatus = Object.values(associated_destinations);
    } catch (e) {
      // silently error
    }
  }

  @action
  async undelete() {
    const { backend, path } = this.args;
    try {
      await this.api.secrets.kvV2UndeleteVersions(path, backend, { versions: [this.version] });
      this.flashMessages.success(`Successfully undeleted ${path}.`);
      this.transition();
    } catch (err) {
      const { message } = await this.api.parseError(err);
      this.flashMessages.danger(`There was a problem undeleting ${path}. Error: ${message}.`);
    }
  }

  @action
  async handleDestruction(type) {
    const { backend, path } = this.args;
    try {
      if (type === 'destroy') {
        await this.api.secrets.kvV2DestroyVersions(path, backend, { versions: [this.version] });
      } else if (type === 'delete-latest-version') {
        await this.api.secrets.kvV2Delete(path, backend);
      } else if (type === 'delete-version') {
        await this.api.secrets.kvV2DeleteVersions(path, backend, { versions: [this.version] });
      } else {
        throw 'type must be one of delete-latest-version, delete-version, or destroy.';
      }
      const verb = type.includes('delete') ? 'deleted' : 'destroyed';
      this.flashMessages.success(`Successfully ${verb} Version ${this.version} of ${path}.`);
      this.transition();
    } catch (err) {
      const { message } = await this.api.parseError(err);
      const verb = type.includes('delete') ? 'deleting' : 'destroying';
      this.flashMessages.danger(
        `There was a problem ${verb} Version ${this.version} of ${path}. Error: ${message}.`
      );
    }
  }

  transition() {
    // transition to the overview to prevent automatically reading sensitive secret data
    this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.index');
  }

  get sortedVersions() {
    return sortedVersions(this.args.metadata?.versions);
  }

  get version() {
    return (
      this.args.secret?.version ||
      this.router.currentRoute.queryParams?.version ||
      this.sortedVersions[0]?.version
    );
  }

  get hideHeaders() {
    return this.showJsonView || this.emptyState;
  }

  get secretState() {
    const { destroyed, created_time } = this.args.secret;
    if (destroyed) return 'destroyed';
    if (this.isSecretDeleted) return 'deleted';
    if (created_time) return 'created';
    return '';
  }

  get versionState() {
    const { secret } = this.args;
    if (secret.failReadErrorCode !== 403) {
      return this.secretState;
    }
    // If the user can't read secret data, get the current version
    // state from metadata versions
    if (this.sortedVersions) {
      const version = this.version;
      const meta = version ? this.sortedVersions.find((v) => v.version == version) : this.sortedVersions[0];
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
    const { canUndelete } = this.args.capabilities;
    if (canUndelete) {
      return this.versionState === 'deleted';
    }
    return false;
  }

  get showDelete() {
    const { canDeleteVersion, canDeleteLatestVersion } = this.args.capabilities;
    if (canDeleteVersion || canDeleteLatestVersion) {
      return this.versionState === 'created' || this.versionState === '';
    }
    return false;
  }

  get isSecretDeleted() {
    return isDeleted(this.args.secret.deletion_time);
  }

  get showDestroy() {
    const { canDestroyVersion } = this.args.capabilities;
    if (canDestroyVersion) {
      return this.versionState !== 'destroyed' && this.version;
    }
    return false;
  }

  get emptyState() {
    const { canReadData, canReadMetadata } = this.args.capabilities;

    if (!canReadData) {
      return {
        title: 'You do not have permission to read this secret',
        message:
          'Your policies may permit you to write a new version of this secret, but do not allow you to read its current contents.',
      };
    }
    // only destructure if we can read secret data
    const { version, destroyed } = this.args.secret;
    if (destroyed) {
      return {
        title: `Version ${version} of this secret has been permanently destroyed`,
        message: `A version that has been permanently deleted cannot be restored. ${
          canReadMetadata
            ? ' You can view other versions of this secret in the Version History tab above.'
            : ''
        }`,
        link: '/vault/docs/secrets/kv/kv-v2',
      };
    }
    if (this.isSecretDeleted) {
      return {
        title: `Version ${version} of this secret has been deleted`,
        message: `This version has been deleted but can be undeleted. ${
          canReadMetadata
            ? 'View other versions of this secret by clicking the Version History tab above.'
            : ''
        }`,
        link: '/vault/docs/secrets/kv/kv-v2',
      };
    }
    return false;
  }
}
