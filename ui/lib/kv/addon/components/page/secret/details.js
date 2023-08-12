/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { next } from '@ember/runloop';

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
  @tracked showJsonView = false;

  @action
  toggleJsonView() {
    this.showJsonView = !this.showJsonView;
  }

  @action
  onClose(dropdown) {
    // strange issue where closing dropdown triggers full transition (which redirects to auth screen in production)
    // closing dropdown in next tick of run loop fixes it
    next(() => dropdown.actions.close());
  }

  @action
  async undelete() {
    const { secret } = this.args;
    try {
      await secret.destroyRecord({
        adapterOptions: { deleteType: 'undelete-version', deleteVersions: secret.version },
      });
      this.flashMessages.success(`Successfully undeleted ${secret.path}`);
      this.router.transitionTo('vault.cluster.secrets.backend.kv.secret', {
        queryParams: { version: secret.version },
      });
    } catch (err) {
      this.flashMessages.danger(`There was a problem un-deleting ${secret.path}`);
    }
  }

  @action
  async handleDestruction(type) {
    const { secret } = this.args;
    // TODO can use secret.state getter for verb once branch is updated
    try {
      await secret.destroyRecord({ adapterOptions: { deleteType: type, deleteVersions: secret.version } });
      const verb = type.includes('delete') ? 'deleted' : 'destroyed';
      this.flashMessages.success(`Successfully ${verb} Version ${secret.version} of ${secret.path}.`);
      this.router.transitionTo('vault.cluster.secrets.backend.kv.secret', {
        queryParams: { version: secret.version },
      });
    } catch (err) {
      const verb = type.includes('delete') ? 'deleting' : 'destroying';
      this.flashMessages.danger(`There was a problem ${verb} Version ${secret.version} of ${secret.path}`);
    }
  }

  get isDeactivated() {
    // TODO can use secret.state getter for verb once branch is updated
    const secret = this.args;
    return !secret.deletionTime && !secret.destroyed;
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
