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
import { verbToPastTense } from 'core/helpers/verb-to-tense';

/**
 * @module KvDeleteModal displays a button for a delete type and launches a modal (undelete is the only mode that does not launch the modal).
 *  All destroyRecord network requests are handled inside the component.
 *
 * <KvDeleteModal
 *  @mode="destroy"
 *  @path="my-secret"
 *  @secret={{this.model.secret}}
 *  @metadata={{this.model.metadata}}
 * />
 *
 * @param {string} mode - Either: delete, destroy, undelete, or delete metadata.
 * @param {string} path - The name of the secret.
 * @param {object} secret - The kv/data model.
 * @param {object} [metadata] - The kv/metadata model. It is only required when mode is "delete" or "metadata-delete".
 */

export default class KvDeleteModal extends Component {
  @service flashMessages;
  @service router;
  @service store;
  @tracked deleteType;
  @tracked modalOpen = false;

  get modalIntro() {
    switch (this.args.mode) {
      case 'delete':
        return 'There are two ways to delete a version of a secret. Both delete actions can be un-deleted later. How would you like to proceed?';
      case 'destroy':
        return `This action will permanently destroy Version ${this.args.secret.version}
        of the secret, and the secret data cannot be read or recovered later.`;
      case 'delete metadata':
        return 'This will permanently delete the metadata and versions of the secret. All version history will be removed. This cannot be undone.';
      default:
        return assert('mode must be one of undelete, delete, destroy, delete metadata.');
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
          ? `The latest version of the secret is already ${this.args.metadata.secretState.state}.`
          : 'You do not have permission to delete the latest version of this secret.',
      },
    ];
  }

  @action handleButtonClick(mode) {
    // undelete is the only use case were a modal is not shown.
    if (mode === 'undelete') {
      this.deleteType = 'undelete-version';
      this.save.perform();
      return;
    }
    this.modalOpen = true;
    // deleteType is the param used in the case/switch for deleteRecord on the adapter.
    if (mode === 'destroy') {
      this.deleteType = 'destroy-version';
    }
    if (mode === 'delete metadata') {
      this.deleteType = 'delete-metadata';
    }
    // if mode is delete, they still need to select what kind of delete operation they'd like to perform so we don't set the deleteType until they select a radio option.
  }

  flashMessageMessage(isSuccess) {
    if (isSuccess) {
      if (this.deleteType === 'delete-metadata') {
        return `Successfully deleted the metadata and all version data for the secret ${this.args.path}.`;
      } else {
        return `Successfully ${verbToPastTense(this.args.mode, 'past')} Version ${
          this.args.secret.version
        } of ${this.args.path}.`;
      }
    } else {
      if (this.deleteType === 'delete-metadata') {
        return `There was an issue deleting ${this.args.path} metadata.`;
      } else {
        return `There was an issue ${verbToPastTense(this.args.mode, 'gerund')} Version ${
          this.args.secret.version
        } of ${this.args.path}.`;
      }
    }
  }

  @(task(function* () {
    this.modalOpen = false;
    try {
      yield this.args.secret.destroyRecord({
        adapterOptions: { deleteType: this.deleteType, deleteVersions: this.args.secret.version },
      });
      this.flashMessages.success(this.flashMessageMessage(true));
      if (this.deleteType === 'delete-metadata') {
        return this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
      }
      return this.router.transitionTo('vault.cluster.secrets.backend.kv.secret', {
        queryParams: { version: this.args.secret.version },
      });
    } catch (err) {
      this.flashMessages.danger(`${this.flashMessageMessage(false)} Error: ${err.message}`);
    }
  }).drop())
  save;
}
