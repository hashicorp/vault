/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { next } from '@ember/runloop';
import { inject as service } from '@ember/service';

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
  @tracked deleteModalOpen = false;
  @tracked destroyModalOpen = false;
  @tracked deleteType = ''; // when set this is either delete-latest-version or delete-version.
  @service flashMessages;
  @service router;

  @action
  toggleJsonView() {
    this.showJsonView = !this.showJsonView;
  }

  @action
  toggleModal(btn) {
    btn === 'delete'
      ? (this.deleteModalOpen = !this.deleteModalOpen)
      : (this.destroyModalOpen = !this.destroyModalOpen);
  }

  @action
  onClose(dropdown) {
    // strange issue where closing dropdown triggers full transition (which redirects to auth screen in production)
    // closing dropdown in next tick of run loop fixes it
    next(() => dropdown.actions.close());
  }

  @action async handleDelete(mode) {
    // mode is either: delete, destroy or undelete.
    const deleteType = mode === 'delete' ? this.deleteType : mode;
    try {
      await this.args.secret.destroyRecord({
        adapterOptions: { deleteType, deleteVersions: this.args.secret.version },
      });
      this.flashMessages.success(
        `Successfully ${this.verbToTense(mode, 'past')} Version ${this.args.secret.version} of ${
          this.args.path
        }.`
      );
      return this.router.transitionTo('vault.cluster.secrets.backend.kv.secret', {
        queryParams: { version: this.args.secret.version },
      });
    } catch {
      return `There was an issue ${this.verbToTense(mode, 'gerund')} Version ${this.args.secret.version} of ${
        this.args.path
      }.`;
    }
  }

  get generateDeleteRadioOptions() {
    let isDeactivated = false;
    if (this.args.secret.canReadMetadata) {
      isDeactivated = this.args.metadata?.currentSecret.isDeactivated;
    }

    return [
      {
        key: 'delete-version',
        label: 'Delete this version',
        description: `This deletes Version ${this.args.secret.version} of the secret.`,
        disabled: !this.args.secret.canDeleteVersion,
        tooltipMessage: 'You do not have permission to delete a specific version.',
      },
      {
        key: 'delete-latest-version',
        label: 'Delete latest version',
        description: 'This deletes the most recent version of the secret.',
        disabled: !this.args.secret.canDeleteLatestVersion || isDeactivated,
        tooltipMessage: isDeactivated
          ? `The latest version of the secret is already ${this.args.metadata.currentSecret.state}.`
          : 'You do not have permission to delete the latest version of this secret.',
      },
    ];
  }

  verbToTense(verb, tense) {
    if (tense === 'gerund') {
      // ending in 'ing ex: delete => deleting || destroy => destroying
      return (
        verb
          .replace(/([^aeiouy])y$/, '$1i')
          .replace(/([^aeiouy][aeiou])([^aeiouy])$/, '$1$2$2')
          .replace(/e$/, '') + 'ing'
      );
    }
    if (tense === 'past') {
      // ending in 'ing ex: delete => deleted || destroy => destroyed
      return (
        verb
          .replace(/([^aeiouy])y$/, '$1i')
          .replace(/([^aeiouy][aeiou])([^aeiouy])$/, '$1$2$2')
          .replace(/e$/, '') + 'ed'
      );
    }
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
