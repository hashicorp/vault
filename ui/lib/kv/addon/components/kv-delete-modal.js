/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';
import { assert } from '@ember/debug';

/**
 * @module KvDeleteModal displays a button for a delete type and launches a the respective modal. All destroyRecord network requests are handled inside the component.
 *
 * <KvDeleteModal
 *  @mode="destroy"
 *  @secret={{this.model.secret}}
 * />
 *
 * @param {string} mode - Either: delete, destroy, or metadata-delete.
 * @param {number} version - This is the version the user is currently viewing.
 * @param {number} currentVersion - This is the current version of the secret.
 * @param {string} path - The name of the secret.
 * @param {boolean} canDeleteVersion - True if the user has UPDATE on the kv/delete endpoint and access to send the version.
 * @param {boolean} canDeleteLatestVersion - True if the user has DELETE on the kv/data endpoint.
 */

export default class KvDeleteModal extends Component {
  @service flashMessages;
  @tracked deleteType;
  @tracked modalOpen = false;

  get modalIntro() {
    switch (this.args.mode) {
      case 'delete':
        return 'There are two ways to delete a version of a secret. Both delete actions can be un-deleted later if you need. How would you like to proceed?';
      case 'destroy':
        return `This action will permanently destroy Version ${this.args.version}
        of the secret, and the secret data cannot be read or recovered later.`;
      case 'metadata-delete':
        return 'This will permanently delete the metadata and versions of the secret. All version history will be removed. This cannot be undone.';
      default:
        return assert('mode must be one of delete, destroy, metadata-delete.');
    }
  }

  get generateRadioDeleteOptions() {
    return [
      {
        key: 'delete-specific-version',
        label: 'Delete this version',
        description: `This deletes Version ${this.args.version} of the secret.`,
        disabled: !this.args.canDeleteVersion,
        tooltipMessage: 'You do not have permission to delete a specific version.',
      },
      {
        key: 'delete-latest-version',
        label: 'Delete latest version',
        description: 'This deletes the most recent version of the secret.',
        disabled: !this.args.canDeleteLatestVersion,
        tooltipMessage: 'You do not have permission to delete the latest version.',
      },
    ];
  }

  @action handleButtonClick(mode) {
    this.modalOpen = true;
    // if mode is destroy, the deleteType is destroy.
    // if mode is delete, they still need to select what kind of delete operation they'd like to perform.
    this.deleteType = mode === 'destroy' ? 'destroy' : '';
  }

  @(task(function* () {
    try {
      yield this.args.secret.destroyRecord({
        adapterOptions: { deleteType: this.deleteType, deleteVersions: [this.args.version] },
      });
      this.flashMessages.success(
        `Successfully ${this.args.mode === 'delete' ? 'deleted' : 'destroyed'} Version ${
          this.args.version
        } of secret ${this.args.path}.`
      );
      this.router.transitionTo('vault.cluster.secrets.backend.kv.secret');
    } catch (err) {
      this.error = err.message;
      this.modalOpen = false;
    }
  }).drop())
  save;
}
