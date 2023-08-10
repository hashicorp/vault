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
 *  @path="my-secret"
 *  @secret={{this.model.secret}}
 * />
 *
 * @param {string} mode - Either: delete, destroy, undelete, or destroy-everything.
 * @param {string} path - The name of the secret.
 * @param {object} secret - The kv/data model.
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
        return 'There are two ways to delete a version of a secret. Both delete actions can be un-deleted later if you need. How would you like to proceed?';
      case 'destroy':
        return `This action will permanently destroy Version ${this.args.secret.version}
        of the secret, and the secret data cannot be read or recovered later.`;
      case 'metadata-delete':
        return 'This will permanently delete the metadata and versions of the secret. All version history will be removed. This cannot be undone.';
      default:
        return assert('mode must be one of undelete, delete, destroy, metadata-delete.');
    }
  }

  get generateRadioDeleteOptions() {
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
        disabled: !this.args.secret.canDeleteLatestVersion,
        tooltipMessage: 'You do not have permission to delete the latest version.',
      },
    ];
  }

  get flashMessageVerb() {
    switch (this.args.mode) {
      case 'delete':
        return 'deleted';
      case 'destroy':
        return 'destroyed';
      case 'undelete':
        return 'undeleted';
      default:
        return assert('mode must be one of undelete, delete, destroy, metadata-delete.');
    }
  }

  @action handleButtonClick(mode) {
    // undelete is the only use case were a modal is not shown.
    if (mode === 'undelete') {
      this.deleteType = 'undelete-version';
      this.save.perform();
      return;
    }
    this.modalOpen = true;
    // if mode is destroy, the deleteType is destroy-version.
    // if mode is delete, they still need to select what kind of delete operation they'd like to perform.
    this.deleteType = mode === 'destroy' ? 'destroy-version' : '';
  }

  @(task(function* () {
    try {
      yield this.args.secret.destroyRecord({
        adapterOptions: { deleteType: this.deleteType, deleteVersions: this.args.secret.version },
      });
      this.flashMessages.success(
        `Successfully ${this.flashMessageVerb} Version ${this.args.secret.version} of ${this.args.path}.`
      );
      this.modalOpen = false;
      this.router.transitionTo('vault.cluster.secrets.backend.kv.secret');
    } catch (err) {
      // ARG TODO handle in flash message.
      this.modalOpen = false;
    }
  }).drop())
  save;
}
