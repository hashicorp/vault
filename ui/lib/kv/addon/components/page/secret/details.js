/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { next } from '@ember/runloop';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module KvSecretDetails renders the key/value data of a KV secret.
 * It also renders a dropdown to display different versions of the secret.
 * <Page::Secret::Details
 *  @path={{this.model.path}}
 *  @secret={{this.model.secret}}
 *  @metadata={{this.model.metadata}}
 *  @breadcrumbs={{this.breadcrumbs}}
  />
 *
 * @param {string} path - path of kv secret 'my/secret' used as the title for the KV page header
 * @param {model} secret - Ember data model: 'kv/data'
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 */

export default class KvSecretDetails extends Component {
  @service flashMessages;
  @service router;
  @service store;

  @tracked showJsonView = false;
  @tracked wrappedData = null;

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

  @action
  async undelete() {
    const { secret } = this.args;
    try {
      await secret.destroyRecord({
        adapterOptions: { deleteType: 'undelete', deleteVersions: this.version },
      });
      this.flashMessages.success(`Successfully undeleted ${secret.path}.`);
      this.refreshRoute();
    } catch (err) {
      this.flashMessages.danger(
        `There was a problem undeleting ${secret.path}. Error: ${err.errors.join(' ')}.`
      );
    }
  }

  @action
  async handleDestruction(type) {
    const { secret } = this.args;
    try {
      await secret.destroyRecord({ adapterOptions: { deleteType: type, deleteVersions: this.version } });
      this.flashMessages.success(`Successfully ${secret.state} Version ${this.version} of ${secret.path}.`);
      this.refreshRoute();
    } catch (err) {
      const verb = type.includes('delete') ? 'deleting' : 'destroying';
      this.flashMessages.danger(
        `There was a problem ${verb} Version ${this.version} of ${secret.path}. Error: ${err.errors.join(
          ' '
        )}.`
      );
    }
  }

  refreshRoute() {
    // transition to the parent secret route to refresh both metadata and data models
    this.router.transitionTo('vault.cluster.secrets.backend.kv.secret', {
      queryParams: { version: this.version },
    });
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
      if (meta?.deletion_time) {
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
    if (!this.args.secret.canReadData) {
      return {
        title: 'You do not have permission to read this secret',
        message:
          'Your policies may permit you to write a new version of this secret, but do not allow you to read its current contents.',
      };
    }
    // only destructure if we can read secret data
    const { version, destroyed, deletionTime } = this.args.secret;
    if (destroyed) {
      return {
        title: `Version ${version} of this secret has been permanently destroyed`,
        message: `A version that has been permanently deleted cannot be restored. ${
          this.args.secret.canReadMetadata
            ? ' You can view other versions of this secret in the Version History tab above.'
            : ''
        }`,
        link: '/vault/docs/secrets/kv/kv-v2',
      };
    }
    if (deletionTime) {
      return {
        title: `Version ${version} of this secret has been deleted`,
        message: `This version has been deleted but can be undeleted. ${
          this.args.secret.canReadMetadata
            ? 'View other versions of this secret by clicking the Version History tab above.'
            : ''
        }`,
        link: '/vault/docs/secrets/kv/kv-v2',
      };
    }
    return false;
  }
}
